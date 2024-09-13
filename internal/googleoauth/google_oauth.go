package googleoauth

import (
	"bytes"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"nine-dubz/internal/file"
	"nine-dubz/internal/user"
	"os"
	"path/filepath"
	"strings"
)

type UseCase struct {
	GoogleOAuthInteractor Interactor
	UserUseCase           *user.UseCase
	FileUseCase           *file.UseCase
}

func New(db *gorm.DB, ur *user.UseCase, fuc *file.UseCase) *UseCase {
	clientId, ok := os.LookupEnv("GOOGLE_CLIENT_ID")
	if !ok {
		fmt.Println("Google Client ID not found in environment")
	}
	clientSecret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET")
	if !ok {
		fmt.Println("Google client secret not found in environment")
	}
	siteUrl, ok := os.LookupEnv("SITE_URL")
	if !ok {
		log.Println("SITE_URL not found in environment")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  siteUrl + "/api/authorize/google/",
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
		FileUseCase: fuc,
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

func (uc *UseCase) Register(registrationRequest *UserRegistrationRequest) (uint, error) {
	registrationPayload := NewUserRegistrationRequest(registrationRequest)
	registrationPayload.Active = true

	resp, err := http.Get(registrationRequest.PictureUrl)
	if err == nil {
		contentType := resp.Header.Get("Content-Type")

		pictureExt := strings.Split(contentType, "/")
		if len(pictureExt) == 2 {
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				bodyReader := bytes.NewReader(body)
				picture, err := uc.FileUseCase.Create(
					bodyReader,
					"google_img."+pictureExt[1],
					filepath.Join("user/google", registrationRequest.Id),
					"public",
				)
				if err == nil {
					registrationPayload.Picture = picture
				}
			}
		}
	}

	return uc.UserUseCase.Add(registrationPayload)
}
