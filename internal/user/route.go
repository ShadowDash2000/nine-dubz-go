package user

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Route("/get-short", func(r chi.Router) {
			r.
				With(h.IsAuthorized).
				Get("/", h.GetUserShortHandler)
		})

		r.Route("/check-by-name", func(r chi.Router) {
			r.Get("/", h.CheckUserWithNameExistsHandler)
		})

		r.Route("/update-picture", func(r chi.Router) {
			r.
				With(h.IsAuthorized).
				Post("/", h.UpdatePictureHandler)
		})
	})
	r.
		With(h.IsNotAuthorized).
		Route("/authorize/inner", func(r chi.Router) {
			r.Route("/register", func(r chi.Router) {
				r.Post("/", h.RegisterHandler)
			})
			r.Route("/login", func(r chi.Router) {
				r.Post("/", h.LoginHandler)
			})
			r.Route("/confirm", func(r chi.Router) {
				r.Get("/", h.ConfirmRegistrationHandler)
			})
		})
}
