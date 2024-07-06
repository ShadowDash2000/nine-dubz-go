package controller

import (
	"fmt"
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

func (rp *RoleRepository) CheckRoutePermission(userName string, routePattern string, method string) (bool, error) {
	roles := &[]model.Role{}

	if len(userName) == 0 {
		result := rp.DB.Preload("ApiMethods", "path = ? AND method = ?", routePattern, method).
			Find(&roles, "code = ?", "all")
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		user := &model.User{}
		result := rp.DB.Preload("Roles.ApiMethods", "path = ? AND method = ?", routePattern, method).
			First(&user, "name = ?", userName)
		if result.Error != nil {
			return false, result.Error
		}
		roles = &user.Roles
	}

	isUserHavePermission := false
	for _, role := range *roles {
		fmt.Println("Role Code: " + role.Code)
		if role.Code == "admin" {
			isUserHavePermission = true
			break
		}

		for _, apiMethod := range role.ApiMethods {
			fmt.Println("API Method: " + apiMethod.Path + ", Method: " + apiMethod.Method)
			isUserHavePermission = true
		}
	}

	return isUserHavePermission, nil
}
