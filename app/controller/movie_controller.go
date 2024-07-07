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

type MovieController struct {
	MovieInteractor usecase.MovieInteractor
}

func NewMovieController(db *gorm.DB) *MovieController {
	return &MovieController{
		MovieInteractor: usecase.MovieInteractor{
			MovieRepository: &MovieRepository{
				DB: db,
			},
		},
	}
}

func (mc *MovieController) Add(w http.ResponseWriter, r *http.Request) {
	movie := &model.Movie{}
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := mc.MovieInteractor.Add(movie)
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 400,
			StatusText:     "Cannot add movie",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}

	render.JSON(w, r, id)
}

func (mc *MovieController) Get(w http.ResponseWriter, r *http.Request) {
	movieId, err := strconv.ParseUint(chi.URLParam(r, "movieId"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid movie id", http.StatusBadRequest)
	}

	movie, err := mc.MovieInteractor.Get(uint(movieId))
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 404,
			StatusText:     "Movie not found",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}

	render.JSON(w, r, movie)
}

func (mc *MovieController) GetAll(w http.ResponseWriter, r *http.Request) {
	pagination := r.Context().Value("pagination").(*model.Pagination)

	movies, err := mc.MovieInteractor.GetAll(pagination)
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 404,
			StatusText:     "Movies not found",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}

	render.JSON(w, r, movies)
}
