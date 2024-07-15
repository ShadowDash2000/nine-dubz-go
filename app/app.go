package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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
	"path/filepath"
	"strings"
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
	// Use cases
	muc := mail.New()
	fuc := file.New(app.DB)
	tuc := token.New(app.DB)
	ruc := role.New(app.DB)
	movuc := movie.New(app.DB)
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
	ta := tokenauthorize.New(tokenSecretKey, "nine-dubz")

	// Http handlers
	uh := user.NewHandler(uuc, tuc, ta, lang)
	fh := file.NewHandler(fuc)
	mh := movie.NewHandler(movuc, uh, fuc)
	goah := googleoauth.NewHandler(goauc, uh, tuc, ta)

	//app.Router.Use(middleware.Logger)
	app.Router.Use(middleware.Recoverer)
	app.Router.Use(middleware.URLFormat)
	app.Router.Use(httprate.Limit(
		10,
		2*time.Second,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
		}),
	))
	app.Router.Use(render.SetContentType(render.ContentTypeJSON))
	app.Router.Use(lang.SetLanguageContext)

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

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "upload/thumbs"))
	FileServer(app.Router, "/upload", filesDir)

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

	fmt.Println(fmt.Sprintf("Starting server at: %s:%s", appIp, appPort))

	err = http.ListenAndServe(appIp+":"+appPort, app.Router)
	if err != nil {
		return
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
