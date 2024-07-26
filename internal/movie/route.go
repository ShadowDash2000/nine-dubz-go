package movie

import (
	"github.com/go-chi/chi/v5"
	"nine-dubz/internal/pagination"
)

func (h *Handler) Routes(r chi.Router) {
	r.Route("/movie", func(r chi.Router) {
		r.
			With(pagination.SetPaginationContextMiddleware).
			Get("/", h.GetMultipleHandler)

		r.
			With(h.UserHandler.IsAuthorized).
			With(h.UserHandler.UserPermission).
			Route("/user", func(r chi.Router) {
				r.Route("/", func(r chi.Router) {
					r.
						With(pagination.SetPaginationContextMiddleware).
						Get("/", h.GetMultipleForUserHandler)
					r.Post("/", h.AddHandler)
				})
				r.Route("/{movieCode}", func(r chi.Router) {
					r.Delete("/", h.DeleteHandler)
					r.Post("/", h.UpdateHandler)
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
	})
}
