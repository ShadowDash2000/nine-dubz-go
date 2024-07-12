package usecase

import (
	"nine-dubz/app/model"
)

type UserRepository interface {
	Add(user *model.User) error
	Remove(id uint) error
	Save(user *model.User) error
	Updates(user *model.User) error
	Get(user *model.User) error
	GetWhere(user *model.User, where map[string]interface{}) error
	GetById(id uint) (*model.User, error)
	GetByName(name string) (*model.User, error)
	Login(user *model.User) bool
	LoginWOPassword(user *model.User) bool
}
