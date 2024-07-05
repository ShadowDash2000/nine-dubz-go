package usecase

import "nine-dubz/app/model"

type RoleRepository interface {
	Add(role *model.Role) (uint, error)
	Remove(id uint) error
	Update(role *model.Role) error
	Get(id uint) (*model.Role, error)
	CheckRoutePermission(userId uint, routePattern string, method string) (bool, error)
}
