package payload

import (
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model"
)

type RegistrationPayload struct {
	*gorm.Model
	*model.User
	PictureId omit `json:"-"`
	Picture   omit `json:"-"`
	Roles     omit `json:"-"`
}

func NewRegistrationPayload(user *model.User) *RegistrationPayload {
	return &RegistrationPayload{
		User: &model.User{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		},
	}
}

func (rp *RegistrationPayload) Bind(r *http.Request) error {
	return nil
}
