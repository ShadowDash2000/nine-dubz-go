package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name       string      `json:"name"`
	Code       string      `json:"code" gorm:"unique"`
	ApiMethods []ApiMethod `json:"apiMethods" gorm:"many2many:role_api_methods;"`
}
