package usecase

import "nine-dubz/app/model/payload"

type HelperRepository interface {
	ValidateUserName(userName string) bool
	ValidateEmail(email string) bool
	ValidatePassword(password string) bool
	ValidateRegistrationFields(user *payload.RegistrationPayload) error
}
