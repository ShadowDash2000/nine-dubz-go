package controller

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model"
	"nine-dubz/app/usecase"
	"strconv"
)

type RoleController struct {
	RoleInteractor  usecase.RoleInteractor
	TokenController *TokenController
}

func NewRoleController(db *gorm.DB, tokenController *TokenController) *RoleController {
	return &RoleController{
		RoleInteractor: usecase.RoleInteractor{
			RoleRepository: &RoleRepository{
				DB: db,
			},
		},
		TokenController: tokenController,
	}
}

func (rc *RoleController) Add(w http.ResponseWriter, r *http.Request) {
	role := &model.Role{}
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := rc.RoleInteractor.Add(role)
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 400,
			StatusText:     "Cannot add role",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}

	render.JSON(w, r, id)
}

func (rc *RoleController) Get(w http.ResponseWriter, r *http.Request) {
	roleId, err := strconv.ParseUint(chi.URLParam(r, "roleId"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid role id", http.StatusBadRequest)
	}

	role, err := rc.RoleInteractor.Get(uint(roleId))
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 404,
			StatusText:     "Role not found",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}

	render.JSON(w, r, role)
}

func (rc *RoleController) Permission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err, http.StatusUnauthorized, "No token cookie"))
			return
		}

		if err = rc.TokenController.Verify(tokenCookie.Value); err != nil {
			render.Render(w, r, ErrInvalidRequest(err, http.StatusUnauthorized, "Can't verify token"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "token", tokenCookie.Value)

		routePattern := chi.RouteContext(r.Context()).RoutePattern()
		method := r.Method

		isUserHavePermission, user, err := rc.RoleInteractor.CheckUserPermission(tokenCookie.Value, routePattern, method)
		if err != nil || !isUserHavePermission {
			render.Render(w, r, ErrInvalidRequest(err, http.StatusForbidden, "You don't have permission"))
			return
		}

		ctx = context.WithValue(ctx, "userId", user.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
