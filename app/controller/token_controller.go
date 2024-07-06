package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"nine-dubz/app/usecase"
)

type TokenController struct {
	TokenInteractor usecase.TokenInteractor
}

func NewTokenController(secret string) *TokenController {
	return &TokenController{
		TokenInteractor: usecase.TokenInteractor{
			TokenRepository: &TokenRepository{
				secret: secret,
			},
		},
	}
}

func (tc *TokenController) Create(userName string) (string, error) {
	return tc.TokenInteractor.Create(userName)
}

func (tc *TokenController) Verify(signedString string) (*jwt.Token, error) {
	return tc.TokenInteractor.Verify(signedString)
}
