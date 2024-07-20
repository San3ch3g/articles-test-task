package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

func GetTokenClaimsFromJWT(tokenString string, secret []byte) (*AuthorClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthorClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthorClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
