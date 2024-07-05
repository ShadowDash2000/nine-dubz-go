package usecase

import "nine-dubz/app/model"

type UserInteractor struct {
	UserRepository UserRepository
}

func (ui *UserInteractor) Add(user *model.User) (uint, error) {
	return ui.UserRepository.Add(user)
}

func (ui *UserInteractor) Remove(id uint) error {
	return ui.UserRepository.Remove(id)
}

func (ui *UserInteractor) Update(user *model.User) error {
	return ui.UserRepository.Update(user)
}

func (ui *UserInteractor) Get(id uint) (*model.User, error) {
	return ui.UserRepository.Get(id)
}
