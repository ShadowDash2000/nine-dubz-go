package app

import (
	"fmt"
	"github.com/alitto/pond"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"log"
	"net/http"
	"nine-dubz/internal/comment"
	"nine-dubz/internal/file"
	"nine-dubz/internal/googleoauth"
	"nine-dubz/internal/mail"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/public"
	"nine-dubz/internal/response"
	"nine-dubz/internal/role"
	"nine-dubz/internal/seo"
	"nine-dubz/internal/subscription"
	"nine-dubz/internal/token"
	"nine-dubz/internal/user"
	"nine-dubz/internal/video"
	"nine-dubz/internal/view"
	"nine-dubz/pkg/etag"
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
	pool := pond.New(4, 300)
	defer pool.StopAndWait()

	// Use cases
	muc := mail.New()
	fuc := file.New(app.DB)
	tuc := token.New(app.DB)
	ruc := role.New(app.DB)
	vuc := view.New(app.DB)
	viduc := video.New(app.DB, fuc)
	uuc := user.New(app.DB, tuc, ruc, fuc, muc)
	subuc := subscription.New(app.DB)
	movuc := movie.New(app.DB, pool, viduc, fuc, vuc, subuc)
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
	ph := public.NewHandler(seouc)
	uh := user.NewHandler(uuc, tuc, ta)
	fh := file.NewHandler(fuc)
	mh := movie.NewHandler(movuc, uh, fuc, ta, tuc)
	goah := googleoauth.NewHandler(goauc, uh, tuc, ta)
	ch := comment.NewHandler(cuc, uh)
	seoh := seo.NewHandler(seouc)
	subh := subscription.NewHandler(subuc, uh)

	//app.Router.Use(middleware.Logger)
	app.Router.Use(
		middleware.Recoverer,
		middleware.URLFormat,
		render.SetContentType(render.ContentTypeJSON),
		language.SetLanguageContext,
		etag.Etag,
	)

	isDev, ok := os.LookupEnv("IS_DEV")
	if !ok {
		isDev = "false"
	}

	app.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if isDev == "true" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	ph.Routes(app.Router)

	app.Router.Route("/api", func(r chi.Router) {
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			response.RenderError(w, r, http.StatusNotFound, "not found")
		})

		fh.Routes(r)

		r.
			With(httprate.Limit(
				30,
				2*time.Second,
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Too many requests", http.StatusTooManyRequests)
				}),
			)).Route("/", func(r chi.Router) {
			uh.Routes(r)
			mh.Routes(r)
			goah.Routes(r)
			ch.Routes(r)
			seoh.Routes(r)
			subh.Routes(r)
		})
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
	go movuc.RetryVideoPostProcess()

	err := http.ListenAndServe(appIp+":"+appPort, app.Router)
	if err != nil {
		return
	}
}
