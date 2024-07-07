package usecase

import (
	"github.com/golang-jwt/jwt/v5"
	"nine-dubz/app/model"
)

type TokenRepository interface {
	Create(user *model.User) (string, *jwt.RegisteredClaims, error)
	Verify(token string) (*jwt.Token, error)
}
