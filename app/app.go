package app

import (
	"fmt"
	"github.com/alitto/pond"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"log"
	"net/http"
	"nine-dubz/internal/comment"
	"nine-dubz/internal/file"
	"nine-dubz/internal/googleoauth"
	"nine-dubz/internal/mail"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/response"
	"nine-dubz/internal/role"
	"nine-dubz/internal/seo"
	"nine-dubz/internal/token"
	"nine-dubz/internal/user"
	"nine-dubz/pkg/language"
	"nine-dubz/pkg/tokenauthorize"
	"os"
	"time"
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
	pool := pond.New(10, 300)
	defer pool.StopAndWait()

	// Use cases
	muc := mail.New()
	fuc := file.New(app.DB)
	tuc := token.New(app.DB)
	ruc := role.New(app.DB)
	movuc := movie.New(app.DB, pool, fuc)
	uuc := user.New(app.DB, tuc, ruc, fuc, muc)
	goauc := googleoauth.New(app.DB, uuc, fuc)
	cuc := comment.New(app.DB, movuc, uuc)
	seouc := seo.New(movuc)

	// JWT Token
	tokenSecretKey, ok := os.LookupEnv("TOKEN_SECRET_KEY")
	if !ok {
		log.Println("TOKEN_SECRET_KEY environment variable not set")
	}
	ta := tokenauthorize.New(tokenSecretKey, "nine-dubz")

	// Http handlers
	uh := user.NewHandler(uuc, tuc, ta)
	fh := file.NewHandler(fuc)
	mh := movie.NewHandler(movuc, uh, fuc, ta, tuc)
	goah := googleoauth.NewHandler(goauc, uh, tuc, ta)
	ch := comment.NewHandler(cuc, uh)
	seoh := seo.NewHandler(seouc)

	//app.Router.Use(middleware.Logger)
	app.Router.Use(middleware.Recoverer)
	app.Router.Use(middleware.URLFormat)
	app.Router.Use(httprate.Limit(
		30,
		2*time.Second,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
		}),
	))
	app.Router.Use(render.SetContentType(render.ContentTypeJSON))
	app.Router.Use(language.SetLanguageContext)

	app.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/**
			TODO remove in future
			*/
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	app.Router.Route("/assets", func(r chi.Router) {
		fs := http.FileServer(http.Dir("public/dist/assets"))
		r.Method("GET", "/*", http.StripPrefix("/assets/", fs))
	})

	app.Router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("public/dist/index.html")
		if err != nil {
			response.RenderError(w, r, http.StatusInternalServerError, "")
			return
		}

		document, err := html.Parse(file)
		if err != nil {
			response.RenderError(w, r, http.StatusInternalServerError, "")
			return
		}

		seouc.SetSeo(r, document)

		html.Render(w, document)
	})

	app.Router.Route("/api", func(r chi.Router) {
		uh.Routes(r)
		fh.Routes(r)
		mh.Routes(r)
		goah.Routes(r)
		ch.Routes(r)
		seoh.Routes(r)
	})

	appIp, ok := os.LookupEnv("APP_IP")
	if !ok {
		appIp = "localhost"
	}
	appPort, ok := os.LookupEnv("APP_PORT")
	if !ok {
		appPort = "8080"
	}

	fmt.Println(fmt.Sprintf("Starting server at: %s:%s", appIp, appPort))

	// If server were crashed, try to re-post-process them
	movuc.RetryVideoPostProcess()

	err := http.ListenAndServe(appIp+":"+appPort, app.Router)
	if err != nil {
		return
	}
}
