package googleoauth

import (
	"gorm.io/gorm"
	"nine-dubz/internal/file"
	"nine-dubz/internal/user"
)

type AuthorizeState struct {
	gorm.Model
	State string `json:"-"`
}

type GoogleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type GoogleToken struct {
	IssuedTo      string `json:"issued_to"`
	Audience      string `json:"audience"`
	UserId        string `json:"user_id"`
	Scope         string `json:"scope"`
	ExpiresIn     int    `json:"expires_in"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	AccessType    string `json:"access_type"`
}

type UserLoginRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserLoginResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

func NewUserLoginRequest(userLoginRequest *UserLoginRequest) *user.User {
	return &user.User{
		Active: true,
		Email:  userLoginRequest.Email,
	}
}

type UserRegistrationRequest struct {
	Id         string
	Name       string
	Email      string
	PictureUrl string
	Picture    *file.File
}

type UserRegistrationResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

func NewUserRegistrationRequest(userRegistrationRequest *UserRegistrationRequest) *user.User {
	return &user.User{
		Name:    userRegistrationRequest.Name,
		Email:   userRegistrationRequest.Email,
		Picture: userRegistrationRequest.Picture,
	}
}
