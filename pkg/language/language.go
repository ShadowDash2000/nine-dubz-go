package language

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
)

type Repository struct {
	Languages []Language
}

type Language struct {
	Code    string    `json:"code"`
	Strings []Message `json:"strings"`
}

type Message struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

func New(languagePath string) (*Repository, error) {
	entries, err := os.ReadDir(languagePath)
	if err != nil {
		return nil, errors.New("no language directory")
	}

	var languages []Language

	for _, e := range entries {
		fileContent, err := os.Open(languagePath + "/" + e.Name())
		if err != nil {
			return nil, err
		}

		buff, _ := io.ReadAll(fileContent)

		language := &Language{}
		if err = json.Unmarshal(buff, &language); err != nil {
			return nil, err
		}

		languages = append(languages, *language)

		fileContent.Close()
	}

	return &Repository{languages}, nil
}

func (rp *Repository) GetLanguageCode(r *http.Request) string {
	return r.Context().Value("lang").(string)
}

func (rp *Repository) GetStringByCode(code string, languageCode string) (string, error) {
	languageIndex := slices.IndexFunc(rp.Languages, func(l Language) bool {
		return l.Code == languageCode
	})

	if languageIndex == -1 {
		return "", errors.New("language not found")
	}

	stringIndex := slices.IndexFunc(rp.Languages[languageIndex].Strings, func(s Message) bool {
		return s.Code == code
	})

	if stringIndex == -1 {
		return "", errors.New("string not found")
	}

	return rp.Languages[languageIndex].Strings[stringIndex].Text, nil
}

func (rp *Repository) GetFormattedStringByCode(code string, values map[string]string, languageCode string) (string, error) {
	languageString, err := rp.GetStringByCode(code, languageCode)
	if err != nil {
		return "", err
	}

	for key, value := range values {
		languageString = strings.ReplaceAll(languageString, "{"+key+"}", value)
	}

	return languageString, nil
}

func (rp *Repository) SetLanguageContext(next http.Handler) http.Handler {
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
