package user

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"nine-dubz/internal/response"
)

func (h *Handler) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			response.RenderError(w, r, http.StatusUnauthorized, "Token cookie not found")
			return
		}

		if _, err = h.TokenAuthorize.VerifyToken(tokenCookie.Value); err != nil {
			response.RenderError(w, r, http.StatusUnauthorized, "Invalid token")
			return
		}

		userId, err := h.TokenUseCase.GetUserIdByToken(tokenCookie.Value)
		if err != nil {
			response.RenderError(w, r, http.StatusUnauthorized, "User not found")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", *(userId))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) TryToGetUserId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId *uint

		tokenCookie, err := r.Cookie("token")
		if err == nil {
			if _, err = h.TokenAuthorize.VerifyToken(tokenCookie.Value); err == nil {
				userId, _ = h.TokenUseCase.GetUserIdByToken(tokenCookie.Value)
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) IsNotAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err == nil {
			if _, err = h.TokenAuthorize.VerifyToken(tokenCookie.Value); err == nil {
				response.RenderError(w, r, http.StatusUnauthorized, "You're authorized already")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) UserPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(uint)

		routePattern := chi.RouteContext(r.Context()).RoutePattern()
		method := r.Method

		if ok := h.UserUseCase.CheckUserPermission(userId, routePattern, method); !ok {
			response.RenderError(w, r, http.StatusForbidden, "Permission denied")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
