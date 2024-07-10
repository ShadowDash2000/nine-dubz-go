package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `json:"name" gorm:"unique;not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"password"`
	PictureId uint   `json:"-"`
	Picture   *File  `json:"picture" gorm:"foreignKey:PictureId;references:ID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Roles     []Role `json:"-" gorm:"many2many:user_roles"`
}
