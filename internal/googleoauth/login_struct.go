package googleoauth

import (
	"nine-dubz/internal/user"
)

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
