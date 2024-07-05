package controller

import (
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
	RoleInteractor usecase.RoleInteractor
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{
		RoleInteractor: usecase.RoleInteractor{
			RoleRepository: &RoleRepository{
				DB: db,
			},
		},
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
