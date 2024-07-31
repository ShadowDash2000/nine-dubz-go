package user

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Route("/get-short", func(r chi.Router) {
			r.
				With(h.IsAuthorized).
				Get("/", h.GetUserShortHandler)

			r.Route("/{userId}", func(r chi.Router) {
				r.
					With(h.IsAuthorized).
					Get("/", h.GetUserHandler)
			})
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
		Route("/authorize/inner", func(r chi.Router) {
			r.
				With(h.IsNotAuthorized).
				Route("/register", func(r chi.Router) {
					r.Post("/", h.RegisterHandler)
				})
			r.
				With(h.IsNotAuthorized).
				Route("/login", func(r chi.Router) {
					r.Post("/", h.LoginHandler)
				})
			r.
				With(h.IsAuthorized).
				Route("/logout", func(r chi.Router) {
					r.Get("/", h.LogoutHandler)
				})
			r.
				With(h.IsNotAuthorized).
				Route("/confirm", func(r chi.Router) {
					r.Get("/", h.ConfirmRegistrationHandler)
				})
		})
}
