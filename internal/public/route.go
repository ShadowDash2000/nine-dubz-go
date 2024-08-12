package public

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) Routes(r chi.Router) {
	r.Route("/assets", func(r chi.Router) {
		r.Get("/*", h.AssetsHandler)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		h.IndexHandler(w, r)
	})

	r.Get("/", h.IndexHandler)
	r.Route("/movie/{movieCode}", func(r chi.Router) {
		r.Get("/", h.IndexHandler)
	})
	r.Route("/studio", func(r chi.Router) {
		r.Get("/", h.IndexHandler)
		r.Route("/edit/{movieCode}", func(r chi.Router) {
			r.Get("/", h.IndexHandler)
		})
	})
}
