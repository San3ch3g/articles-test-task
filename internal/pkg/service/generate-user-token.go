package service

import (
	"articleModule/internal/pkg/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AuthorClaims struct {
	AuthorId uint32 `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateUserToken(secret []byte, author models.Author) (string, error) {
	claims := AuthorClaims{
		AuthorId: author.Id,
		Username: author.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
