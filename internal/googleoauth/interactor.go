package googleoauth

type Interactor interface {
	GetConsentPageUrl() string
	Authorize(code string, state string) (*GoogleUserInfo, error)
}
