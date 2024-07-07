package payload

import (
	"net/http"
	"nine-dubz/app/model"
)

type UserShortPayload struct {
	Password omit       `json:"password,omitempty"`
	Picture  model.File `json:"picture,omitempty"`
	*UserPayload
}

func NewUserShortPayload(user *model.User) *UserShortPayload {
	return &UserShortPayload{
		UserPayload: &UserPayload{
			User: &model.User{
				Name:    user.Name,
				Email:   user.Email,
				Picture: user.Picture,
			},
		},
	}
}

func (rp *UserShortPayload) Bind(r *http.Request) error {
	return nil
}
