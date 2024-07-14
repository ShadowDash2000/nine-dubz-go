package user

import (
	"gorm.io/gorm"
	"nine-dubz/internal/file"
	"nine-dubz/internal/role"
)

type User struct {
	gorm.Model
	Active    bool `gorm:"default:false"`
	ID        uint
	Name      string      `json:"name" gorm:"unique;not null"`
	Email     string      `json:"email" gorm:"unique;not null"`
	Password  string      `json:"password"`
	Hash      string      `json:"-"`
	PictureId *uint       `json:"-"`
	Picture   *file.File  `json:"picture" gorm:"foreignKey:PictureId;references:ID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Roles     []role.Role `json:"-" gorm:"many2many:user_roles"`
}
