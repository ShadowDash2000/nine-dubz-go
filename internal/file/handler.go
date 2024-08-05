package file

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"nine-dubz/internal/response"
	"strconv"
)

type Handler struct {
	FileUseCase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{
		FileUseCase: uc,
	}
}
func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")

	buff, err := h.FileUseCase.GetFile(fileName)
	if err != nil {
		response.RenderError(w, r, http.StatusNotFound, "No such file")
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(buff)))
	w.Write(buff)
}
