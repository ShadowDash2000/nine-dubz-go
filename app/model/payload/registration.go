package payload

import (
	"nine-dubz/app/model"
)

func NewRegistrationPayload(user *model.User) *model.User {
	return &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}
