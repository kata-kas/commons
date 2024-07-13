package commons

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("invalid token signature")
		}
		return "", errors.New("invalid token")
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Email, nil
}
