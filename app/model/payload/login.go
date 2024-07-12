package payload

import (
	"nine-dubz/app/model"
)

func NewLoginPayload(user *model.User) *model.User {
	return &model.User{
		Active:   user.Active,
		Email:    user.Email,
		Password: user.Password,
	}
}
