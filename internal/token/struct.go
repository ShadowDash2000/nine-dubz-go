package token

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	Token  string `json:"token"`
	UserId uint   `json:"userId"`
}
