package usecase

import (
	"nine-dubz/app/model/payload"
)

type GoogleOauthInteractor struct {
	GoogleOauthRepository GoogleOauthRepository
}

func (goi *GoogleOauthInteractor) GetConsentPageUrl() string {
	return goi.GoogleOauthRepository.GetConsentPageUrl()
}

func (goi *GoogleOauthInteractor) Authorize(code string, state string) (*payload.GoogleUserInfo, error) {
	return goi.GoogleOauthRepository.Authorize(code, state)
}
