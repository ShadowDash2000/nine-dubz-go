package usecase

import "nine-dubz/app/model"

type UserRepository interface {
	Add(user *model.User) (uint, error)
	Remove(id uint) error
	Update(user *model.User) error
	Get(id uint) (*model.User, error)
}
