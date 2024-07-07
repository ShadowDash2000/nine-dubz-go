package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"nine-dubz/app/model"
	"nine-dubz/app/usecase"
)

type TokenController struct {
	TokenInteractor usecase.TokenInteractor
	DB              *gorm.DB
}

func NewTokenController(secret string, db *gorm.DB) *TokenController {
	return &TokenController{
		TokenInteractor: usecase.TokenInteractor{
			TokenRepository: &TokenRepository{
				secret: secret,
				DB:     db,
			},
		},
	}
}

func (tc *TokenController) Create(user *model.User) (string, *jwt.RegisteredClaims, error) {
	signedString, claims, err := tc.TokenInteractor.Create(user)

	return signedString, claims, err
}

func (tc *TokenController) Verify(token string) error {
	_, err := tc.TokenInteractor.Verify(token)

	return err
}
