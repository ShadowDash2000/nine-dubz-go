package usecase

type GoogleOauthInteractor struct {
	GoogleOauthRepository GoogleOauthRepository
}

func (goi *GoogleOauthInteractor) GetConsentPageUrl() string {
	return goi.GoogleOauthRepository.GetConsentPageUrl()
}

func (goi *GoogleOauthInteractor) Authorize(code string) error {
	return goi.GoogleOauthRepository.Authorize(code)
}
