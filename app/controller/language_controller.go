package controller

import (
	"context"
	"net/http"
	"nine-dubz/app/usecase"
)

type LanguageController struct {
	LanguageInteractor usecase.LanguageInteractor
}

func NewLanguageController(languagePath string) *LanguageController {
	return &LanguageController{
		LanguageInteractor: usecase.LanguageInteractor{
			LanguageRepository: NewLanguageRepository(languagePath),
		},
	}
}

func (lc *LanguageController) GetLanguageCode(r *http.Request) string {
	return r.Context().Value("lang").(string)
}

func (lc *LanguageController) GetStringByCode(r *http.Request, code string) (string, error) {
	return lc.LanguageInteractor.GetStringByCode(code, lc.GetLanguageCode(r))
}

func (lc *LanguageController) Language(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		languageCode := ""
		languageCookie, err := r.Cookie("lang")
		if err != nil {
			languageCode = "eng"
		} else {
			languageCode = languageCookie.Value
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "lang", languageCode)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
