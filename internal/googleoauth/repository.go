package googleoauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

type Repository struct {
	DB          *gorm.DB
	OauthConfig *oauth2.Config
}

func (gor *Repository) GetConsentPageUrl() string {
	buff := make([]byte, 10)
	io.ReadFull(rand.Reader, buff)
	state := base64.StdEncoding.EncodeToString(buff)

	result := gor.DB.Create(&AuthorizeState{
		State: state,
	})
	if result.Error != nil {
		return ""
	}

	return gor.OauthConfig.AuthCodeURL(state)
}

func (gor *Repository) Authorize(code string, state string) (*GoogleUserInfo, error) {
	authorizeState := &AuthorizeState{}
	maxStateLifeTime := time.Now().Add(-2 * time.Hour)
	result := gor.DB.Where("created_at > ?", maxStateLifeTime).First(authorizeState, "state = ?", state)
	if result.Error != nil {
		return nil, result.Error
	}

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

	userInfo := &GoogleUserInfo{}
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
