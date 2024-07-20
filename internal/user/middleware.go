package user

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Token cookie not found", http.StatusUnauthorized)
			return
		}

		if _, err = h.TokenAuthorize.VerifyToken(tokenCookie.Value); err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userId, err := h.TokenUseCase.GetUserIdByToken(tokenCookie.Value)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) TryToGetUSerId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if _, err = h.TokenAuthorize.VerifyToken(tokenCookie.Value); err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userId, err := h.TokenUseCase.GetUserIdByToken(tokenCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
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
				http.Error(w, "You're authorized already", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) UserPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId")
		if userId == "" {
			return
		}

		routePattern := chi.RouteContext(r.Context()).RoutePattern()
		method := r.Method

		if ok := h.UserUseCase.CheckUserPermission(userId.(uint), routePattern, method); !ok {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
