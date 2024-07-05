package controller

import (
	"gorm.io/gorm"
	"nine-dubz/app/model"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		DB: db,
	}
}

func (rp *RoleRepository) Add(role *model.Role) (uint, error) {
	result := rp.DB.Create(role)

	return role.ID, result.Error
}

func (rp *RoleRepository) Remove(id uint) error {
	result := rp.DB.Delete(&model.Role{}, id)

	return result.Error
}

func (rp *RoleRepository) Update(role *model.Role) error {
	result := rp.DB.Save(role)

	return result.Error
}

func (rp *RoleRepository) Get(id uint) (*model.Role, error) {
	role := &model.Role{}
	result := rp.DB.First(role, id)

	return role, result.Error
}

func (rp *RoleRepository) CheckRoutePermission(id uint, routePattern string, method string) (bool, error) {
	user := &model.User{}
	result := rp.DB.Preload("Roles.ApiMethods", "path = ? AND method = ?", routePattern, method).
		First(&user, id)
	if result.Error != nil {
		return false, result.Error
	}

	isUserHavePermission := false
	for _, role := range user.Roles {
		if role.Code == "admin" {
			isUserHavePermission = true
			break
		}

		for range role.ApiMethods {
			isUserHavePermission = true
		}
	}

	return isUserHavePermission, result.Error
}
