package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model"
	"nine-dubz/app/model/payload"
	"nine-dubz/app/usecase"
	"strconv"
)

type MovieController struct {
	MovieInteractor usecase.MovieInteractor
	FileController  *FileController
}

func NewMovieController(db *gorm.DB, fc *FileController) *MovieController {
	return &MovieController{
		MovieInteractor: usecase.MovieInteractor{
			MovieRepository: &MovieRepository{
				DB: db,
			},
		},
		FileController: fc,
	}
}

func (mc *MovieController) AddHandler(w http.ResponseWriter, r *http.Request) {
	movie := &model.Movie{}
	err := mc.MovieInteractor.Add(movie)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "Cannot add movie"))
		return
	}

	file, err := mc.FileController.SocketVideoUpload(w, r)
	if err != nil {
		mc.MovieInteractor.Remove(movie.ID)

		fmt.Println(err)
		return
	}

	movie.Video = file
	mc.MovieInteractor.Updates(movie)
}

func (mc *MovieController) GetHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := strconv.ParseUint(chi.URLParam(r, "movieId"), 10, 32)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, "No movie id value"))
		return
	}

	movie, err := mc.MovieInteractor.Get(uint(movieId))
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err, http.StatusNotFound, "Movie not found"))
		return
	}

	render.JSON(w, r, payload.NewMoviePayload(movie))
}

func (mc *MovieController) GetAllHandler(w http.ResponseWriter, r *http.Request) {
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

	if len(*movies) == 0 {
		render.JSON(w, r, struct{}{})
		return
	}

	var moviesPayload []*model.Movie
	for _, movie := range *movies {
		moviesPayload = append(moviesPayload, payload.NewMoviePayload(&movie))
	}

	render.JSON(w, r, moviesPayload)
}
