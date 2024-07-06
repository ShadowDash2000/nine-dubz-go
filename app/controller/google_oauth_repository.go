package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"nine-dubz/app/model"
)

type GoogleOauthRepository struct {
	OauthConfig *oauth2.Config
}

func (gor *GoogleOauthRepository) GetConsentPageUrl() string {
	return gor.OauthConfig.AuthCodeURL("state-token")
}

func (gor *GoogleOauthRepository) Authorize(code string) error {
	fmt.Println("CODE: " + code)
	token, err := gor.OauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	fmt.Println("TOKEN: " + token.AccessToken)
	request, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Bearer "+token.AccessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	userInfo := &model.GoogleUserInfo{}
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return err
	}

	fmt.Println(userInfo)

	return nil
}
