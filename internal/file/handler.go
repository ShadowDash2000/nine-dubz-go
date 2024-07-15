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

func (h *Handler) StreamFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	requestRange := r.Header.Get("Range")

	buff, contentRange, contentLength, err := h.FileUseCase.StreamFile(fileName, requestRange)
	if err != nil {
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}

	w.Header().Set("Accept-Ranges", "bytes")
	if len(requestRange) > 0 {
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("Content-Length", contentLength)
	}
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusPartialContent)
	w.Write(buff)
}
