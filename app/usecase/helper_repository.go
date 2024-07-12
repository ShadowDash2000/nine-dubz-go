package usecase

import (
	"nine-dubz/app/model"
)

type HelperRepository interface {
	ValidateUserName(userName string) bool
	ValidateEmail(email string) bool
	ValidatePassword(password string) bool
	ValidateRegistrationFields(user *model.User) error
}
