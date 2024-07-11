package usecase

import (
	"nine-dubz/app/model"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func (ui *UserInteractor) Add(user *model.User) (uint, error) {
	return ui.UserRepository.Add(user)
}

func (ui *UserInteractor) Remove(id uint) error {
	return ui.UserRepository.Remove(id)
}

func (ui *UserInteractor) Save(user *model.User) error {
	return ui.UserRepository.Save(user)
}

func (ui *UserInteractor) Updates(user *model.User) error {
	return ui.UserRepository.Updates(user)
}

func (ui *UserInteractor) Get(id uint) (*model.User, error) {
	return ui.UserRepository.Get(id)
}

func (ui *UserInteractor) GetByName(name string) (*model.User, error) {
	return ui.UserRepository.GetByName(name)
}

func (ui *UserInteractor) Login(user *model.User) bool {
	return ui.UserRepository.Login(user)
}

func (ui *UserInteractor) LoginWOPassword(user *model.User) bool {
	return ui.UserRepository.LoginWOPassword(user)
}
