package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"log"
	"net/http"
	"nine-dubz/internal/file"
	"nine-dubz/internal/googleoauth"
	"nine-dubz/internal/mail"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/role"
	"nine-dubz/internal/token"
	"nine-dubz/internal/user"
	"nine-dubz/pkg/language"
	"nine-dubz/pkg/tokenauthorize"
	"os"
)

type App struct {
	Router *chi.Mux
	DB     *gorm.DB
}

func NewApp(db gorm.DB) *App {
	return &App{
		Router: chi.NewRouter(),
		DB:     &db,
	}
}

func (app *App) Start() {
	// Use cases
	muc := mail.New()
	fuc := file.New(app.DB)
	tuc := token.New(app.DB)
	ruc := role.New(app.DB)
	movuc := movie.New(app.DB, fuc)
	uuc := user.New(app.DB, tuc, ruc, fuc, muc)
	goauc := googleoauth.New(app.DB, uuc)

	// Language translations
	lang, err := language.New("lang")
	if err != nil {
		log.Println(err)
	}

	// JWT Token
	tokenSecretKey, ok := os.LookupEnv("TOKEN_SECRET_KEY")
	if !ok {
		log.Println("TOKEN_SECRET_KEY environment variable not set")
	}
	ta := tokenauthorize.New(tokenSecretKey)

	// Http handlers
	uh := user.NewHandler(uuc, tuc, ta, lang)
	fh := file.NewHandler(fuc)
	mh := movie.NewHandler(movuc, ta, uh, fh)
	goah := googleoauth.NewHandler(goauc, tuc, ta)

	//app.Router.Use(middleware.Logger)
	app.Router.Use(middleware.Recoverer)
	app.Router.Use(middleware.URLFormat)
	app.Router.Use(render.SetContentType(render.ContentTypeJSON))
	app.Router.Use(lang.SetLanguageContextMiddleware)

	app.Router.Use(func(next http.Handler) http.Handler {
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

	app.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
	app.Router.Get("/socket-test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/socket_test.html")
	})

	app.Router.Route("/api", func(r chi.Router) {
		uh.Routes(r)
		fh.Routes(r)
		mh.MovieRoutes(r)
		goah.Routes(r)
	})

	appIp, ok := os.LookupEnv("APP_IP")
	if !ok {
		appIp = "localhost"
	}
	appPort, ok := os.LookupEnv("APP_PORT")
	if !ok {
		appPort = "8080"
	}

	err = http.ListenAndServe(appIp+":"+appPort, app.Router)
	if err != nil {
		return
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`         // user-level status message
	AppCode    int64  `json:"code,omitempty"` // application-specific error code
	ErrorText  string `json:"-"`              // application-level error message, for debugging
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
