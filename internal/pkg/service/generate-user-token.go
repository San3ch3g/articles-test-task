package service

import (
	"articleModule/internal/pkg/models"
	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	AuthorId uint32 `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateUserToken(secret []byte, author models.Author) (string, error) {
	claims := UserClaims{
		AuthorId: author.Id,
		Username: author.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
