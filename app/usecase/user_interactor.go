package usecase

import (
	"nine-dubz/app/model"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func (ui *UserInteractor) Add(user *model.User) error {
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

func (ui *UserInteractor) Get(user *model.User) error {
	return ui.UserRepository.Get(user)
}

func (ui *UserInteractor) GetWhere(user *model.User, where map[string]interface{}) error {
	return ui.UserRepository.GetWhere(user, where)
}

func (ui *UserInteractor) GetById(id uint) (*model.User, error) {
	return ui.UserRepository.GetById(id)
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
