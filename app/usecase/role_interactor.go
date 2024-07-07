package usecase

import "nine-dubz/app/model"

type RoleInteractor struct {
	RoleRepository RoleRepository
}

func (ri *RoleInteractor) Add(role *model.Role) (uint, error) {
	return ri.RoleRepository.Add(role)
}

func (ri *RoleInteractor) Remove(id uint) error {
	return ri.RoleRepository.Remove(id)
}

func (ri *RoleInteractor) Update(role *model.Role) error {
	return ri.RoleRepository.Update(role)
}

func (ri *RoleInteractor) Get(id uint) (*model.Role, error) {
	return ri.RoleRepository.Get(id)
}

func (ri *RoleInteractor) CheckUserPermission(token string, routePattern string, method string) (bool, *model.User, error) {
	return ri.RoleRepository.CheckUserPermission(token, routePattern, method)
}
