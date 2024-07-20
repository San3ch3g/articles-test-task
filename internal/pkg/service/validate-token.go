package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(secret []byte, tokenString string) (*AuthorClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthorClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthorClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
