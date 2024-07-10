package payload

import (
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model"
)

type UserShortPayload struct {
	*gorm.Model
	*model.User
	Password  omit `json:"-"`
	PictureId omit `json:"-"`
	Roles     omit `json:"-"`
}

func NewUserShortPayload(user *model.User) *UserShortPayload {
	return &UserShortPayload{
		User: &model.User{
			Name:    user.Name,
			Email:   user.Email,
			Picture: user.Picture,
		},
	}
}

func (rp *UserShortPayload) Bind(r *http.Request) error {
	return nil
}
