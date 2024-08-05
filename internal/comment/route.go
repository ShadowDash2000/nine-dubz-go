package comment

import (
	"github.com/go-chi/chi/v5"
	"nine-dubz/internal/pagination"
	"nine-dubz/internal/sort"
)

func (h *Handler) Routes(r chi.Router) {
	r.Route("/comment", func(r chi.Router) {
		r.
			Route("/{movieCode}", func(r chi.Router) {
				r.
					With(h.UserHandler.IsAuthorized).
					Post("/", h.AddCommentHandler)
				r.
					With(pagination.SetPaginationContextMiddleware).
					With(sort.SetSortContextMiddleware).
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
