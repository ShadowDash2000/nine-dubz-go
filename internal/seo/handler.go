package seo

import (
	"github.com/go-chi/render"
	"net/http"
	"nine-dubz/internal/response"
)

type Handler struct {
	SeoUseCase *UseCase
}

func NewHandler(seouc *UseCase) *Handler {
	return &Handler{
		SeoUseCase: seouc,
	}
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		response.RenderError(w, r, http.StatusBadRequest, "No path parameter")
		return
	}

	seo, err := h.SeoUseCase.GetSeo(path, r)
	if err != nil {
		response.RenderError(w, r, http.StatusInternalServerError, "")
	}

	render.JSON(w, r, NewGetSeoResponse(seo))
}
