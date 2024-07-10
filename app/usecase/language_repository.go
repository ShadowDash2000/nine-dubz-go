package usecase

type LanguageRepository interface {
	GetStringByCode(code string, languageCode string) (string, error)
}
