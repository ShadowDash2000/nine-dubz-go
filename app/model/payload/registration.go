package payload

import (
	"net/http"
	"nine-dubz/app/model"
)

type RegistrationPayload struct {
	*UserPayload
}

func NewRegistrationPayload(user *model.User) *RegistrationPayload {
	return &RegistrationPayload{
		UserPayload: &UserPayload{
			User: &model.User{
				Name:     user.Name,
				Email:    user.Email,
				Password: user.Password,
			}},
	}
}

func (rp *RegistrationPayload) Bind(r *http.Request) error {
	return nil
}
