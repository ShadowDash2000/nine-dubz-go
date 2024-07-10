package usecase

type LanguageInteractor struct {
	LanguageRepository LanguageRepository
}

func (li *LanguageInteractor) GetStringByCode(code string, languageCode string) (string, error) {
	return li.LanguageRepository.GetStringByCode(code, languageCode)
}
