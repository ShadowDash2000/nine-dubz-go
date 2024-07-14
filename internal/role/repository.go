package role

import (
	"gorm.io/gorm"
	"nine-dubz/internal/apimethod"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetApiMethodsByRolesIds(rolesIds []uint) ([]apimethod.ApiMethod, error) {
	var roles []Role
	result := r.DB.Preload("ApiMethods").Find(&roles, rolesIds)
	if result.Error != nil {
		return nil, result.Error
	}

	var apiMethods []apimethod.ApiMethod
	for _, role := range roles {
		for _, apiMethod := range role.ApiMethods {
			apiMethods = append(apiMethods, apiMethod)
		}
	}

	return apiMethods, nil
}
