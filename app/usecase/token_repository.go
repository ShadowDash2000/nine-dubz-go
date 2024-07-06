package usecase

import "github.com/golang-jwt/jwt/v5"

type TokenRepository interface {
	Create(userName string) (string, error)
	Verify(token string) (*jwt.Token, error)
}
