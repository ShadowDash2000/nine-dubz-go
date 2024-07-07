package controller

import (
	"fmt"
	"github.com/go-chi/render"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"nine-dubz/app/usecase"
	"os"
)

type GoogleOauthController struct {
	GoogleOAuthInteractor usecase.GoogleOauthInteractor
}

func NewGoogleOauthController() *GoogleOauthController {
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
		RedirectURL:  "http://localhost:25565/authorize",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOauthController{
		GoogleOAuthInteractor: usecase.GoogleOauthInteractor{
			GoogleOauthRepository: &GoogleOauthRepository{
				OauthConfig: oauthConfig,
			},
		},
	}
}

func (goc *GoogleOauthController) GetConsentPageUrl(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, struct {
		Url string `json:"url"`
	}{goc.GoogleOAuthInteractor.GetConsentPageUrl()})
}

func (goc *GoogleOauthController) Authorize(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err := goc.GoogleOAuthInteractor.Authorize(code)
	if err != nil {
		errResponse := &ErrResponse{
			Err:            err,
			HTTPStatusCode: 500,
			StatusText:     "",
			ErrorText:      err.Error(),
		}
		render.Render(w, r, errResponse)
		return
	}
}
