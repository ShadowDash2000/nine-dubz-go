package googleoauth

import (
	"nine-dubz/internal/user"
)

type UserRegistrationRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserRegistrationResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

func NewUserRegistrationRequest(userRegistrationRequest *UserRegistrationRequest) *user.User {
	return &user.User{
		Name:  userRegistrationRequest.Name,
		Email: userRegistrationRequest.Email,
	}
}
