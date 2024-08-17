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
	ID      uint       `json:"id"`
	Name    string     `json:"name"`
	Email   string     `json:"email"`
	Picture *file.File `json:"picture"`
}

func NewShortResponse(user *User) *ShortResponse {
	return &ShortResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Picture: user.Picture,
	}
}

type GetPublicResponse struct {
	ID      uint       `json:"id"`
	Name    string     `json:"name"`
	Picture *file.File `json:"picture"`
}

func NewGetPublicResponse(user *User) *GetPublicResponse {
	return &GetPublicResponse{
		ID:      user.ID,
		Name:    user.Name,
		Picture: user.Picture,
	}
}

func NewGetPublicResponseMultiple(users []User) []*GetPublicResponse {
	var result []*GetPublicResponse
	for _, user := range users {
		result = append(result, NewGetPublicResponse(&user))
	}

	return result
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLoginRequest(user *LoginRequest) *User {
	return &User{
		Active:   true,
		Email:    user.Email,
		Password: user.Password,
	}
}

type RegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRegistrationRequest(user *RegistrationRequest) *User {
	return &User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

type UpdateRequest struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
}

func NewUpdateRequest(user *UpdateRequest) *User {
	return &User{
		ID:   user.ID,
		Name: user.Name,
	}
}
