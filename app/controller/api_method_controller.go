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

type ApiMethodController struct {
	ApiMethodInteractor usecase.ApiMethodInteractor
}

func NewApiMethodController(db *gorm.DB) *ApiMethodController {
	return &ApiMethodController{
		ApiMethodInteractor: usecase.ApiMethodInteractor{
			ApiMethodRepository: &ApiMethodRepository{
				DB: db,
			},
		},
	}
}

func (uc *ApiMethodController) AddHandler(w http.ResponseWriter, r *http.Request) {
	apiMethod := &model.ApiMethod{}
	err := json.NewDecoder(r.Body).Decode(&apiMethod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := uc.ApiMethodInteractor.Add(apiMethod)
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 400,
			StatusText:     "Cannot add api method",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}

	render.JSON(w, r, id)
}

func (uc *ApiMethodController) GetHandler(w http.ResponseWriter, r *http.Request) {
	apiMethodId, err := strconv.ParseUint(chi.URLParam(r, "apiMethodId"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid api method id", http.StatusBadRequest)
	}

	apiMethod, err := uc.ApiMethodInteractor.Get(uint(apiMethodId))
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 404,
			StatusText:     "Api method not found",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}

	render.JSON(w, r, apiMethod)
}
