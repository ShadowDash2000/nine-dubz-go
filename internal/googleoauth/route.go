package googleoauth

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes(r chi.Router) {
	r.Route("/authorize/google", func(r chi.Router) {
		r.With(h.UserHandler.IsNotAuthorized).Get("/", h.Authorize)
		r.Get("/get-url", h.GetConsentPageUrlHandler)
	})
}
