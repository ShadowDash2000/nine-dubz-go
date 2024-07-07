package payload

import (
	"net/http"
	"nine-dubz/app/model"
)

type LoginPayload struct {
	Name omit `json:"name,omitempty"`
	*UserPayload
}

func NewLoginPayload(user *model.User) *LoginPayload {
	return &LoginPayload{
		UserPayload: &UserPayload{
			User: &model.User{
				Email:    user.Email,
				Password: user.Password,
			},
		},
	}
}

func (rp *LoginPayload) Bind(r *http.Request) error {
	return nil
}
