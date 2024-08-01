package seo

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes(r chi.Router) {
	r.Route("/seo", func(r chi.Router) {
		r.Get("/", h.GetHandler)
	})
}
