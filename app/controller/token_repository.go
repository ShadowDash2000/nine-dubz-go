package controller

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"nine-dubz/app/model"
	"time"
)

type TokenRepository struct {
	secret string
	DB     *gorm.DB
}

func (tr *TokenRepository) Create(user *model.User) (string, *jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   user.Email,
		Issuer:    "nine-dubz",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(tr.secret))
	if err != nil {
		return "", claims, err
	}

	result := tr.DB.Create(&model.Token{
		Token:  signedString,
		UserId: user.ID,
	})
	if result.Error != nil {
		return "", claims, result.Error
	}

	return signedString, claims, err
}

func (tr *TokenRepository) Verify(signedString string) (*jwt.Token, error) {
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tr.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
