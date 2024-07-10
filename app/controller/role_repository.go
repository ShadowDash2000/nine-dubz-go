package controller

import (
	"gorm.io/gorm"
	"nine-dubz/app/model"
)

type RoleRepository struct {
	DB *gorm.DB
}

func (rp *RoleRepository) Add(role *model.Role) (uint, error) {
	result := rp.DB.Create(role)

	return role.ID, result.Error
}

func (rp *RoleRepository) Remove(id uint) error {
	result := rp.DB.Delete(&model.Role{}, id)

	return result.Error
}

func (rp *RoleRepository) Save(role *model.Role) error {
	result := rp.DB.Save(role)

	return result.Error
}

func (rp *RoleRepository) Get(id uint) (*model.Role, error) {
	role := &model.Role{}
	result := rp.DB.First(role, id)

	return role, result.Error
}

func (rp *RoleRepository) CheckUserPermission(tokenString string, routePattern string, method string) (bool, *model.User) {
	token := &model.Token{}
	result := rp.DB.First(&token, "token = ?", tokenString)
	if result.Error != nil {
		return false, nil
	}

	user := &model.User{}
	result = rp.DB.
		Preload("Roles.ApiMethods", "path = ? AND method = ?", routePattern, method).
		First(&user, token.UserId)
	if result.Error != nil {
		return false, nil
	}
	roles := &user.Roles

	isUserHavePermission := false
	for _, role := range *roles {
		if role.Code == "admin" {
			isUserHavePermission = true
			break
		}

		for range role.ApiMethods {
			isUserHavePermission = true
		}
	}

	return isUserHavePermission, user
}
