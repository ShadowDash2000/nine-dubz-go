package public

import (
	"golang.org/x/net/html"
	"net/http"
	"nine-dubz/internal/response"
	"nine-dubz/internal/seo"
	"os"
	"path/filepath"
)

type Handler struct {
	DistPath   string
	SeoUseCase *seo.UseCase
}

func NewHandler(seouc *seo.UseCase) *Handler {
	distPath, ok := os.LookupEnv("DIST_PATH")
	if !ok {
		distPath = "public/dist"
	}

	return &Handler{
		DistPath:   distPath,
		SeoUseCase: seouc,
	}
}

func (h *Handler) AssetsHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir(filepath.Join(h.DistPath, "assets")))
	fs = http.StripPrefix("/assets/", fs)
	fs.ServeHTTP(w, r)
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(filepath.Join(h.DistPath, "index.html"))
	if err != nil {
		response.RenderError(w, r, http.StatusInternalServerError, "")
		return
	}

	document, err := html.Parse(file)
	if err != nil {
		response.RenderError(w, r, http.StatusInternalServerError, "")
		return
	}

	h.SeoUseCase.SetSeo(r, document)

	html.Render(w, document)
}
