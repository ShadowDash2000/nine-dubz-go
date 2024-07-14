package movie

import (
	"github.com/go-chi/chi/v5"
	"nine-dubz/middleware"
)

func (h *Handler) MovieRoutes(r chi.Router) {
	r.Route("/movie", func(r chi.Router) {
		r.With(middleware.PaginationMiddleware).Get("/", h.GetAllHandler)

		r.Route("/upload", func(r chi.Router) {
			r.
				With(h.TokenAuthorize.IsAuthorizedMiddleware).
				With(h.UserHandler.PermissionMiddleware).
				Get("/", h.AddHandler)
		})

		r.Route("/{movieId}", func(r chi.Router) {
			r.Get("/", h.GetHandler)
		})
	})
}
