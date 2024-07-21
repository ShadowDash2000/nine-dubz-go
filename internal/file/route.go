package file

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) Routes(r chi.Router) {
	r.Route("/file", func(r chi.Router) {
		r.Route("/{fileName}", func(r chi.Router) {
			r.With(CacheControl).Get("/", h.GetFile)
		})
	})
}

func CacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=604800")

		next.ServeHTTP(w, r)
	})
}
