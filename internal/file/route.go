package file

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes(r chi.Router) {
	r.Route("/file", func(r chi.Router) {
		r.Route("/stream/{fileName}", func(r chi.Router) {
			r.Get("/", h.StreamFile)
		})
	})
}
