package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Roles    []Role `json:"-" gorm:"many2many:user_roles"`
}
