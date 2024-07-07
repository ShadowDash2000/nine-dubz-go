package usecase

import (
	"github.com/golang-jwt/jwt/v5"
	"nine-dubz/app/model"
)

type TokenInteractor struct {
	TokenRepository TokenRepository
}

func (ti *TokenInteractor) Create(user *model.User) (string, *jwt.RegisteredClaims, error) {
	return ti.TokenRepository.Create(user)
}

func (ti *TokenInteractor) Verify(token string) (*jwt.Token, error) {
	return ti.TokenRepository.Verify(token)
}
