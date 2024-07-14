package role

import (
	"gorm.io/gorm"
	"nine-dubz/internal/apimethod"
)

type Role struct {
	gorm.Model
	Name       string                `json:"name"`
	Code       string                `json:"code" gorm:"unique"`
	ApiMethods []apimethod.ApiMethod `json:"apiMethods" gorm:"many2many:role_api_methods;"`
}
