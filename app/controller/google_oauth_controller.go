package controller

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model"
	"nine-dubz/app/model/payload"
	"nine-dubz/app/usecase"
	"os"
	"strings"
)

type GoogleOauthController struct {
	GoogleOAuthInteractor usecase.GoogleOauthInteractor
	LanguageController    *LanguageController
	UserController        *UserController
}

func NewGoogleOauthController(db *gorm.DB, lc *LanguageController, uc *UserController) *GoogleOauthController {
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

	return &GoogleOauthController{
		GoogleOAuthInteractor: usecase.GoogleOauthInteractor{
			GoogleOauthRepository: &GoogleOauthRepository{
				OauthConfig: oauthConfig,
				DB:          db,
			},
		},
		LanguageController: lc,
		UserController:     uc,
	}
}

func (goc *GoogleOauthController) GetConsentPageUrlHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, struct {
		Url string `json:"url"`
	}{goc.GoogleOAuthInteractor.GetConsentPageUrl()})
}

func (goc *GoogleOauthController) Authorize(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		text, _ := goc.LanguageController.GetStringByCode(r, "GOOGLE_NO_AUTHORIZE_CODE")
		render.Render(w, r, ErrInvalidRequest(errors.New(""), http.StatusInternalServerError, text))
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		text, _ := goc.LanguageController.GetStringByCode(r, "GOOGLE_NO_AUTHORIZE_STATE")
		render.Render(w, r, ErrInvalidRequest(errors.New(""), http.StatusInternalServerError, text))
		return
	}

	googleUser, err := goc.GoogleOAuthInteractor.Authorize(code, state)
	if err != nil {
		text, _ := goc.LanguageController.GetStringByCode(r, "GOOGLE_CANT_AUTHORIZE")
		render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, text))
		return
	}

	user := &model.User{
		Name:  strings.Split(googleUser.Email, "@")[0],
		Email: googleUser.Email,
	}

	// Try to log in
	loginPayload := payload.NewLoginPayload(user)
	tokenCookie, stringCode, err := goc.UserController.Login(loginPayload, false)
	if err == nil {
		http.SetCookie(w, tokenCookie)
		http.Redirect(w, r, "/", http.StatusOK)
		return
	}

	// Try to register
	registrationPayload := payload.NewRegistrationPayload(user)
	userId, stringCode, err := goc.UserController.Register(registrationPayload, true)
	if err != nil {
		text, _ := goc.LanguageController.GetStringByCode(r, stringCode.Code)
		render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, text))
		return
	} else {
		user.ID = userId
		tokenCookie, err := goc.UserController.GetTokenCookie(user)
		if err != nil {
			text, _ := goc.LanguageController.GetStringByCode(r, "TOKEN_FAILED_TO_CREATE")
			render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, text))
			return
		}

		http.SetCookie(w, tokenCookie)
		http.Redirect(w, r, "/", http.StatusOK)
	}
}
