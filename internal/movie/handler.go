package movie

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"nine-dubz/internal/file"
	"nine-dubz/internal/user"
	"nine-dubz/model"
	"nine-dubz/pkg/tokenauthorize"
	"strconv"
)

type Handler struct {
	MovieUseCase   *UseCase
	TokenAuthorize *tokenauthorize.TokenAuthorize
	UserHandler    *user.Handler
	FileHandler    *file.Handler
}

func NewHandler(uc *UseCase, ta *tokenauthorize.TokenAuthorize, uh *user.Handler, fh *file.Handler) *Handler {
	return &Handler{
		MovieUseCase:   uc,
		TokenAuthorize: ta,
		UserHandler:    uh,
		FileHandler:    fh,
	}
}

func (h *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	movie, err := h.MovieUseCase.Add()
	if err != nil {
		//render.Render(w, r, controller.ErrInvalidRequest(err, http.StatusBadRequest, "Cannot add movie"))
		return
	}

	file, err := h.FileHandler.SocketVideoUpload(w, r)
	if err != nil {
		h.MovieUseCase.Remove(movie.ID)
		return
	}

	movieRequest := &UpdateRequest{
		ID:    movie.ID,
		Video: file,
	}
	h.MovieUseCase.Updates(movieRequest)
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := strconv.ParseUint(chi.URLParam(r, "movieId"), 10, 32)
	if err != nil {
		//render.Render(w, r, controller.ErrInvalidRequest(err, http.StatusBadRequest, "No movie id value"))
		return
	}

	movie, err := h.MovieUseCase.Get(uint(movieId))
	if err != nil {
		//render.Render(w, r, controller.ErrInvalidRequest(err, http.StatusNotFound, "Movie not found"))
		return
	}

	render.JSON(w, r, movie)
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	pagination := r.Context().Value("pagination").(*model.Pagination)

	moviesResponse, err := h.MovieUseCase.GetAll(pagination)
	if err != nil {
		//render.Render(w, r, controller.ErrInvalidRequest(err, http.StatusNotFound, "Movies not found"))
		return
	}

	render.JSON(w, r, moviesResponse)
}
