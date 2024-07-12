package usecase

import (
	"nine-dubz/app/model"
)

type HelperInteractor struct {
	HelperRepository HelperRepository
}

func (hi *HelperInteractor) ValidateUserName(userName string) bool {
	return hi.HelperRepository.ValidateUserName(userName)
}

func (hi *HelperInteractor) ValidateEmail(email string) bool {
	return hi.HelperRepository.ValidateUserName(email)
}

func (hi *HelperInteractor) ValidatePassword(password string) bool {
	return hi.HelperRepository.ValidateUserName(password)
}

func (hi *HelperInteractor) ValidateRegistrationFields(user *model.User) error {
	return hi.HelperRepository.ValidateRegistrationFields(user)
}
