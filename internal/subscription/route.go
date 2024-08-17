package subscription

import (
	"github.com/go-chi/chi/v5"
	"nine-dubz/internal/pagination"
)

func (h *Handler) Routes(r chi.Router) {
	r.
		With(h.UserHandler.IsAuthorized).
		Route("/subscription", func(r chi.Router) {
			r.Route("/{channelId}", func(r chi.Router) {
				r.Post("/", h.SubscribeHandler)
				r.Delete("/", h.UnsubscribeHandler)
			})

			r.
				With(pagination.SetPaginationContextMiddleware).
				Get("/", h.GetMultipleHandler)
		})
}
