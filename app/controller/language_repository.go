package controller

import (
	"encoding/json"
	"errors"
	"io"
	"nine-dubz/app/model/language"
	"os"
	"slices"
)

type LanguageRepository struct {
	LanguagePath string
	Languages    []language.Language
}

func NewLanguageRepository(languagePath string) *LanguageRepository {
	entries, err := os.ReadDir(languagePath)
	if err != nil {
		panic(err)
	}

	var languages []language.Language

	for _, e := range entries {
		fileContent, err := os.Open(languagePath + "/" + e.Name())

		if err != nil {
			panic(err)
		}

		byteResult, _ := io.ReadAll(fileContent)

		languageModel := &language.Language{}
		if err = json.Unmarshal(byteResult, &languageModel); err != nil {
			panic(err)
		}

		languages = append(languages, *languageModel)

		fileContent.Close()
	}

	return &LanguageRepository{
		LanguagePath: languagePath,
		Languages:    languages,
	}
}

func (lr *LanguageRepository) GetStringByCode(code string, languageCode string) (string, error) {
	languageIndex := slices.IndexFunc(lr.Languages, func(l language.Language) bool {
		return l.Code == languageCode
	})

	if languageIndex == -1 {
		return "", errors.New("language not found")
	}

	stringIndex := slices.IndexFunc(lr.Languages[languageIndex].Strings, func(s language.String) bool {
		return s.Code == code
	})

	if stringIndex == -1 {
		return "", errors.New("string not found")
	}

	return lr.Languages[languageIndex].Strings[stringIndex].Text, nil
}
