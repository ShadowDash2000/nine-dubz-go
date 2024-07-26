package user

import "nine-dubz/internal/role"

type Interactor interface {
	Add(user *User) (uint, error)
	Remove(id uint) error
	Save(user *User) error
	Updates(user *User) error
	Get(user *User) error
	GetWhere(user *User, where map[string]interface{}) error
	GetById(id uint) (*User, error)
	GetByName(name string) (*User, error)
	GetRolesByUserId(userId uint) ([]role.Role, error)
	Login(user *User) uint
	LoginWOPassword(user *User) uint
}
