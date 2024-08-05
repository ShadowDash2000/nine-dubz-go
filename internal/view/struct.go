package view

import (
	"gorm.io/gorm"
	"nine-dubz/internal/user"
)

type View struct {
	gorm.Model
	MovieID *uint
	UserID  *uint
	User    user.User
	IP      string
}
