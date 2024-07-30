package comment

import (
	"github.com/go-chi/chi/v5"
	"nine-dubz/internal/pagination"
)

func (h *Handler) Routes(r chi.Router) {
	r.Route("/comment", func(r chi.Router) {
		r.
			With(h.UserHandler.IsAuthorized).
			Route("/{movieCode}", func(r chi.Router) {
				r.
					With(h.UserHandler.IsAuthorized).
					Post("/", h.AddCommentHandler)
				r.
					With(pagination.SetPaginationContextMiddleware).
					With(h.UserHandler.TryToGetUserId).
					Get("/", h.GetMultipleHandler)
				r.
					Route("/{commentId}", func(r chi.Router) {
						r.
							With(h.UserHandler.IsAuthorized).
							Post("/", h.AddCommentHandler)
						r.
							With(pagination.SetPaginationContextMiddleware).
							With(h.UserHandler.TryToGetUserId).
							Get("/", h.GetMultipleSubCommentsHandler)
						r.
							With(h.UserHandler.IsAuthorized).
							Delete("/", h.DeleteCommentHandler)
					})
			})
	})
}