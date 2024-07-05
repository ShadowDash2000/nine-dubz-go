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

type UserController struct {
	UserInteractor usecase.UserInteractor
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		UserInteractor: usecase.UserInteractor{
			UserRepository: &UserRepository{
				DB: db,
			},
		},
	}
}

func (uc *UserController) Add(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := uc.UserInteractor.Add(user)
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 400,
			StatusText:     "Cannot add user",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}

	render.JSON(w, r, id)
}

func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
	}

	user, err := uc.UserInteractor.Get(uint(userId))
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 404,
			StatusText:     "User not found",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}

	render.JSON(w, r, user)
}
