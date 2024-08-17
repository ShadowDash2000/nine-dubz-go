package movie

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"nine-dubz/internal/pagination"
	"nine-dubz/internal/sorting"
)

func (h *Handler) Routes(r chi.Router) {
	r.Route("/movie", func(r chi.Router) {
		r.
			With(pagination.SetPaginationContextMiddleware).
			With(sorting.SetSortContextMiddleware).
			Get("/", h.GetMultipleHandler)

		r.
			With(h.UserHandler.IsAuthorized).
			With(h.UserHandler.UserPermission).
			Route("/user", func(r chi.Router) {
				r.Route("/", func(r chi.Router) {
					r.
						With(pagination.SetPaginationContextMiddleware).
						With(sorting.SetSortContextMiddleware).
						Get("/", h.GetMultipleForUserHandler)
					r.Post("/", h.AddHandler)
				})
				r.Route("/{movieCode}", func(r chi.Router) {
					r.Delete("/", h.DeleteHandler)
					r.
						With(middleware.RequestSize(2<<20)).
						Post("/", h.UpdateHandler)
					r.Get("/", h.GetForUserHandler)
				})
				r.Route("/upload", func(r chi.Router) {
					r.Get("/", h.UploadVideoHandler)
				})
				r.Route("/multiple", func(r chi.Router) {
					r.Delete("/", h.DeleteMultipleHandler)

					r.Route("/status", func(r chi.Router) {
						r.Post("/", h.UpdatePublishStatusHandler)
					})
				})
			})

		r.Route("/{movieCode}", func(r chi.Router) {
			r.
				With(h.UserHandler.TryToGetUserId).
				Get("/", h.GetHandler)
		})
		r.Route("/stream/{movieCode}", func(r chi.Router) {
			r.
				With(h.UserHandler.TryToGetUserId).
				Get("/", h.StreamFile)

			r.
				With(h.UserHandler.TryToGetUserId).
				Head("/", h.StreamFile)
		})
		r.Route("/subscription", func(r chi.Router) {
			r.
				With(pagination.SetPaginationContextMiddleware).
				With(h.UserHandler.IsAuthorized).
				Get("/", h.GetMultipleSubscribedHandler)
		})
	})
}
