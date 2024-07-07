package controller

import (
	"errors"
	"nine-dubz/app/model/payload"
	"regexp"
)

type HelperRepository struct{}

func (hr *HelperRepository) ValidateUserName(userName string) bool {
	matched, err := regexp.MatchString(`^[a-zа-яёA-ZA-ЯЁ0-9]{2,}$`, userName)
	if err != nil || !matched {
		return false
	}
	return true
}

func (hr *HelperRepository) ValidateEmail(email string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`, email)
	if err != nil || !matched {
		return false
	}
	return true
}

func (hr *HelperRepository) ValidatePassword(password string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9$~@#%*!&?=()]{8,}$`, password)
	if err != nil || !matched {
		return false
	}
	return true
}

func (hr *HelperRepository) ValidateRegistrationFields(user *payload.RegistrationPayload) error {
	if !hr.ValidateUserName(user.Name) {
		return errors.New("incorrect user name")
	}
	if !hr.ValidateEmail(user.Email) {
		return errors.New("incorrect email")
	}
	if !hr.ValidatePassword(user.Password) {
		return errors.New("incorrect password")
	}

	return nil
}
