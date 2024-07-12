package controller

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model/payload"
)

type GoogleOauthRepository struct {
	DB          *gorm.DB
	OauthConfig *oauth2.Config
}

func (gor *GoogleOauthRepository) GetConsentPageUrl() string {
	// TODO Make random state
	return gor.OauthConfig.AuthCodeURL("state-token")
}

func (gor *GoogleOauthRepository) Authorize(code string, state string) (*payload.GoogleUserInfo, error) {
	// TODO State check

	token, err := gor.OauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token.AccessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	userInfo := &payload.GoogleUserInfo{}
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
