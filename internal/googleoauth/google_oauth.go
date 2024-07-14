package googleoauth

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"nine-dubz/internal/user"
	"os"
)

type UseCase struct {
	GoogleOAuthInteractor Interactor
	UserUseCase           *user.UseCase
}

func New(db *gorm.DB, ur *user.UseCase) *UseCase {
	clientId, ok := os.LookupEnv("GOOGLE_CLIENT_ID")
	if !ok {
		fmt.Println("Google Client ID not found in environment")
	}
	clientSecret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET")
	if !ok {
		fmt.Println("Google client secret not found in environment")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:25565/api/authorize/google/",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return &UseCase{
		GoogleOAuthInteractor: &Repository{
			DB:          db,
			OauthConfig: oauthConfig,
		},
		UserUseCase: ur,
	}
}

func (uc *UseCase) GetConsentPageUrl() string {
	return uc.GoogleOAuthInteractor.GetConsentPageUrl()
}

func (uc *UseCase) Authorize(code string, state string) (*GoogleUserInfo, error) {
	return uc.GoogleOAuthInteractor.Authorize(code, state)
}

func (uc *UseCase) Login(loginRequest *UserLoginRequest) uint {
	loginPayload := NewUserLoginRequest(loginRequest)

	return uc.UserUseCase.LoginWOPassword(loginPayload)
}

func (uc *UseCase) Register(registrationRequest *UserRegistrationRequest) uint {
	registrationPayload := NewUserRegistrationRequest(registrationRequest)
	registrationPayload.Active = true

	return uc.UserUseCase.Add(registrationPayload)
}
