package usecase

import "github.com/golang-jwt/jwt/v5"

type TokenInteractor struct {
	TokenRepository TokenRepository
}

func (ti *TokenInteractor) Create(userName string) (string, error) {
	return ti.TokenRepository.Create(userName)
}

func (ti *TokenInteractor) Verify(token string) (*jwt.Token, error) {
	return ti.TokenRepository.Verify(token)
}
