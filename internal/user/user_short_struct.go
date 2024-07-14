package user

import (
	"nine-dubz/internal/file"
)

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
