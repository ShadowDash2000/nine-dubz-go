package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"net/http"
)

type RouterController struct {
	Router              chi.Mux
	DB                  *gorm.DB
	MovieController     MovieController
	RoleController      RoleController
	ApiMethodController ApiMethodController
	UserController      UserController
	TokenController     TokenController
}

func NewRouterController(db gorm.DB) *RouterController {
	return &RouterController{
		Router:              *chi.NewRouter(),
		DB:                  &db,
		MovieController:     *NewMovieController(&db),
		RoleController:      *NewRoleController(&db),
		ApiMethodController: *NewApiMethodController(&db),
		UserController:      *NewUserController(&db),
		TokenController:     *NewTokenController("nine-dubz-token-secret"),
	}
}

func (rc *RouterController) HandleRoute() *chi.Mux {
	rc.Router.Use(middleware.Logger)
	rc.Router.Use(middleware.Recoverer)
	rc.Router.Use(middleware.URLFormat)
	rc.Router.Use(render.SetContentType(render.ContentTypeJSON))

	rc.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	rc.Router.Route("/api", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				/**
				TODO remove in future
				*/
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

				next.ServeHTTP(w, r)
			})
		})
		r.Route("/movie", func(r chi.Router) {
			r.With(rc.Permission).Post("/", rc.MovieController.Add)
			r.With(rc.Permission).Get("/", rc.MovieController.GetAll)

			r.Route("/{movieId}", func(r chi.Router) {
				r.With(rc.Permission).Get("/", rc.MovieController.Get)
			})
		})
		r.Route("/role", func(r chi.Router) {
			r.With(rc.Permission).Post("/", rc.RoleController.Add)

			r.Route("/{roleId}", func(r chi.Router) {
				r.With(rc.Permission).Get("/", rc.RoleController.Get)
			})
		})
		r.Route("/api-method", func(r chi.Router) {
			r.With(rc.Permission).Post("/", rc.ApiMethodController.Add)

			r.Route("/{apiMethodId}", func(r chi.Router) {
				r.With(rc.Permission).Get("/", rc.ApiMethodController.Get)
			})
		})
		r.Route("/user", func(r chi.Router) {
			r.With(rc.Permission).Post("/", rc.UserController.Add)

			r.Route("/{userId}", func(r chi.Router) {
				r.With(rc.Permission).Get("/", rc.UserController.Get)
			})
		})
	})

	/*fmt.Println(docgen.MarkdownRoutesDoc(&rc.Router, docgen.MarkdownOpts{
		ProjectPath: "github.com/go-chi/chi/v5",
		Intro:       "Welcome to the chi/_examples/rest generated docs.",
	}))*/

	return &rc.Router
}

func (rc *RouterController) Permission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName := ""
		tokenCookie, err := r.Cookie("token")
		if err == nil {
			token, err := rc.TokenController.Verify(tokenCookie.String())
			if err != nil {
				http.Error(w, http.StatusText(403), 403)
			}
			userName, _ = token.Claims.GetSubject()
		}

		routePattern := chi.RouteContext(r.Context()).RoutePattern()
		method := r.Method

		isUserHavePermission, err := rc.RoleController.RoleInteractor.CheckRoutePermission(userName, routePattern, method)
		if err != nil || !isUserHavePermission {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
