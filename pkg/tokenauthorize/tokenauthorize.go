package tokenauthorize

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type TokenAuthorize struct {
	Issuer    string
	SecretKey string
}

func New(secretKey, issuer string) *TokenAuthorize {
	return &TokenAuthorize{
		Issuer:    issuer,
		SecretKey: secretKey,
	}
}

func (ta *TokenAuthorize) CreateToken(subject string) (string, *jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   subject,
		Issuer:    ta.Issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(ta.SecretKey))
	if err != nil {
		return "", claims, err
	}

	return signedString, claims, err
}

func (ta *TokenAuthorize) VerifyToken(signedString string) (*jwt.Token, error) {
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ta.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (ta *TokenAuthorize) GetTokenCookie(subject string) (*http.Cookie, error) {
	signedString, claims, err := ta.CreateToken(subject)
	if err != nil {
		return nil, err
	}

	tokenCookie := http.Cookie{
		Name:     "token",
		Value:    signedString,
		Path:     "/",
		Expires:  claims.ExpiresAt.Time,
		MaxAge:   0,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteDefaultMode,
	}

	return &tokenCookie, nil
}
