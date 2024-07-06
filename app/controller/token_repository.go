package controller

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenRepository struct {
	secret string
}

func (tr *TokenRepository) Create(userName string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userName,
		Issuer:    "nine-dubz",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(tr.secret)
	if err != nil {
		return "", err
	}

	return signedString, err
}

func (tr *TokenRepository) Verify(signedString string) (*jwt.Token, error) {
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		return tr.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
