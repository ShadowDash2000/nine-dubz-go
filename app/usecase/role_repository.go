package usecase

import "nine-dubz/app/model"

type RoleRepository interface {
	Add(role *model.Role) (uint, error)
	Remove(id uint) error
	Save(role *model.Role) error
	Get(id uint) (*model.Role, error)
	CheckUserPermission(token string, routePattern string, method string) (bool, *model.User)
}
