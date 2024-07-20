package file

import (
	"github.com/go-chi/chi/v5"
	"net/http"
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
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}

	w.Write(buff)
}
