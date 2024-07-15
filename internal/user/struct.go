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

type ShortResponse struct {
	Name    string     `json:"name"`
	Email   string     `json:"email"`
	Picture *file.File `json:"picture"`
}

func NewShortResponse(user *User) *ShortResponse {
	return &ShortResponse{
		Name:    user.Name,
		Email:   user.Email,
		Picture: user.Picture,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

func NewLoginRequest(user *LoginRequest) *User {
	return &User{
		Active:   true,
		Email:    user.Email,
		Password: user.Password,
	}
}

func NewLoginResponse(isSuccess bool) *LoginResponse {
	return &LoginResponse{
		IsSuccess: isSuccess,
	}
}

type RegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

func NewRegistrationRequest(user *RegistrationRequest) *User {
	return &User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

func NewRegistrationResponse(isSuccess bool) *RegistrationResponse {
	return &RegistrationResponse{
		IsSuccess: isSuccess,
	}
}
