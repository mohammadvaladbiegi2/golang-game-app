package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const hmacSampleSecret = "mamad-server"

func BulidToken(name string, id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		"id":   id,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
