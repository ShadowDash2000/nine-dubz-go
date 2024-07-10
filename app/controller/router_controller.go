package controller

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model"
	"strconv"
)

type RouterController struct {
	Router                chi.Mux
	DB                    *gorm.DB
	LanguageController    LanguageController
	MovieController       MovieController
	RoleController        RoleController
	ApiMethodController   ApiMethodController
	UserController        UserController
	TokenController       TokenController
	GoogleOauthController GoogleOauthController
	FileController        FileController
}

func NewRouterController(db gorm.DB) *RouterController {
	tc := *NewTokenController("nine-dubz-token-secret", &db)
	lc := *NewLanguageController("lang")
	mc := *NewMovieController(&db)
	rc := *NewRoleController(&db, &tc)
	amc := *NewApiMethodController(&db)
	uc := *NewUserController(&db, &tc, &lc)
	goc := *NewGoogleOauthController(&db, &lc, &uc)
	fc := *NewFileController(&db, &lc)

	return &RouterController{
		Router:                *chi.NewRouter(),
		DB:                    &db,
		LanguageController:    lc,
		MovieController:       mc,
		RoleController:        rc,
		ApiMethodController:   amc,
		UserController:        uc,
		TokenController:       tc,
		GoogleOauthController: goc,
		FileController:        fc,
	}
}

func (rc *RouterController) HandleRoute() *chi.Mux {
	rc.Router.Use(middleware.Logger)
	rc.Router.Use(middleware.Recoverer)
	rc.Router.Use(middleware.URLFormat)
	rc.Router.Use(render.SetContentType(render.ContentTypeJSON))
	rc.Router.Use(rc.LanguageController.Language)

	rc.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
	rc.Router.Get("/socket-test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/socket_test.html")
	})

	rc.Router.Get("/authorize", rc.GoogleOauthController.Authorize)

	rc.Router.Route("/api", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				/**
				TODO remove in future
				*/
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

				next.ServeHTTP(w, r)
			})
		})

		r.Route("/movie", func(r chi.Router) {
			r.With(rc.RoleController.Permission).Post("/", rc.MovieController.AddHandler)
			r.With(rc.Pagination).Get("/", rc.MovieController.GetAllHandler)

			r.Route("/{movieId}", func(r chi.Router) {
				r.Get("/", rc.MovieController.Get)
			})
		})
		r.Route("/role", func(r chi.Router) {
			r.With(rc.RoleController.Permission).Post("/", rc.RoleController.AddHandler)

			r.Route("/{roleId}", func(r chi.Router) {
				r.With(rc.RoleController.Permission).Get("/", rc.RoleController.GetHandler)
			})
		})
		r.Route("/api-method", func(r chi.Router) {
			r.With(rc.RoleController.Permission).Post("/", rc.ApiMethodController.AddHandler)

			r.Route("/{apiMethodId}", func(r chi.Router) {
				r.With(rc.RoleController.Permission).Get("/", rc.ApiMethodController.GetHandler)
			})
		})

		r.Route("/user", func(r chi.Router) {
			r.With(rc.RoleController.Permission).Post("/", rc.UserController.AddHandler)

			r.Route("/get-short", func(r chi.Router) {
				r.With(rc.RoleController.Permission).Get("/", rc.UserController.GetUserShortHandler)
			})

			r.Route("/get/{userId}", func(r chi.Router) {
				r.With(rc.RoleController.Permission).Get("/", rc.UserController.GetHandler)
			})

			r.Get("/check-by-name", rc.UserController.CheckUserWithNameExistsHandler)
		})

		/**
		TODO Add restriction if user is already authorized
		*/
		r.Route("/authorize", func(r chi.Router) {
			r.Route("/google", func(r chi.Router) {
				r.Get("/", rc.GoogleOauthController.Authorize)
				r.Get("/get-url", rc.GoogleOauthController.GetConsentPageUrlHandler)
			})
			r.Route("/inner", func(r chi.Router) {
				r.Post("/register", rc.UserController.RegisterHandler)
				r.Post("/login", rc.UserController.LoginHandler)
			})
		})
	})

	rc.Router.Route("/file", func(r chi.Router) {
		r.Get("/", rc.FileController.UploadVideoHandler)
	})

	// Generate api doc file
	/*file, err := os.Create("routes.md")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	apiDoc := docgen.MarkdownRoutesDoc(&rc.Router, docgen.MarkdownOpts{
		ProjectPath: "github.com/go-chi/chi/v5",
		Intro:       "Welcome to the chi/_examples/rest generated docs.",
		URLMap: map[string]string{
			"github.com/newline-sandbox/go-chi-docgen-example/vendor/github.com/go-chi/chi/v5/": "https://github.com/go-chi/chi/blob/master/",
		},
		ForceRelativeLinks: true,
	})

	if _, err = file.Write([]byte(apiDoc)); err != nil {
		log.Fatal(err)
	}*/

	return &rc.Router
}

func (rc *RouterController) Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit <= 0 {
			limit = -1
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil || offset < 0 || limit <= 0 {
			offset = -1
		}

		pagination := &model.Pagination{
			Limit:  limit,
			Offset: offset,
		}

		ctx := context.WithValue(r.Context(), "pagination", pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
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

func ErrInvalidRequest(err error, code int, text string) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: code,
		StatusText:     text,
		ErrorText:      err.Error(),
	}
}
