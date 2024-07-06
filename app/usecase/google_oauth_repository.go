package usecase

type GoogleOauthRepository interface {
	GetConsentPageUrl() string
	Authorize(code string) error
}
