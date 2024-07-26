package language

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Language struct {
	Code     string    `json:"code"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

func GetLanguage(languageCode string) (*Language, error) {
	languagePath := GetPath()

	entries, err := os.ReadDir(languagePath)
	if err != nil {
		return nil, errors.New("language: no language directory")
	}

	var language *Language

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		if strings.TrimSuffix(e.Name(), ".json") != languageCode {
			continue
		}

		fileContent, err := os.Open(filepath.Join(languagePath, e.Name()))
		if err != nil {
			return nil, err
		}

		buff, _ := io.ReadAll(fileContent)

		if err = json.Unmarshal(buff, &language); err != nil {
			return nil, err
		}

		fileContent.Close()
	}

	if language == nil {
		return nil, errors.New("language: no language file")
	}

	return language, nil
}

func GetLanguageCode(r *http.Request) string {
	return r.Context().Value("lang").(string)
}

func GetPath() string {
	path, ok := os.LookupEnv("LANG_PATH")
	if !ok {
		return "lang"
	}

	return path
}

func GetMessage(messageCode, languageCode string) (string, error) {
	language, err := GetLanguage(languageCode)
	if err != nil {
		return "", errors.New("language: language not found")
	} else if language.Messages == nil {
		return "", errors.New("language: no language messages")
	}

	messageIndex := slices.IndexFunc(language.Messages, func(m Message) bool {
		return m.Code == messageCode
	})

	if messageIndex == -1 {
		return "", errors.New("language: message not found")
	}

	return language.Messages[messageIndex].Text, nil
}

func GetFormattedMessage(messageCode string, values map[string]string, languageCode string) (string, error) {
	languageString, err := GetMessage(messageCode, languageCode)
	if err != nil {
		return "", err
	}

	for key, value := range values {
		languageString = strings.ReplaceAll(languageString, "{"+key+"}", value)
	}

	return languageString, nil
}

func SetLanguageContext(next http.Handler) http.Handler {
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
