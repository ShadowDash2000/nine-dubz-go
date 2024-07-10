package payload

import (
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model"
)

type LoginPayload struct {
	*gorm.Model
	*model.User
	Name    omit `json:"-"`
	Picture omit `json:"-"`
	Roles   omit `json:"-"`
}

func NewLoginPayload(user *model.User) *LoginPayload {
	return &LoginPayload{
		User: &model.User{
			Email:    user.Email,
			Password: user.Password,
		},
	}
}

func (rp *LoginPayload) Bind(r *http.Request) error {
	return nil
}
