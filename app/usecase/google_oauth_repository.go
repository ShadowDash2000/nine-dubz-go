package usecase

import (
	"nine-dubz/app/model/payload"
)

type GoogleOauthRepository interface {
	GetConsentPageUrl() string
	Authorize(code string, state string) (*payload.GoogleUserInfo, error)
}
