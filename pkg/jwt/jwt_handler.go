package jwt

import (
	"fmt"
	"gamegolang/pkg/richerror"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const hmacSampleSecret = "mamad-server"

type CustomClaims struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
	jwt.RegisteredClaims
}

func BuildToken(name string, id uint) (string, richerror.RichError) {
	claims := CustomClaims{
		Name: name,
		ID:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", richerror.NewError(err, 500, "failed to sign token", map[string]interface{}{"opration": "BuildToken"})
	}

	return tokenString, richerror.RichError{}
}

func VerifyToken(tokenString string) (*CustomClaims, richerror.RichError) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		return nil, richerror.NewError(err, 400, "failed to parse token", map[string]interface{}{"opration": "VerifyToken"})
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, richerror.RichError{}
	}

	return nil, richerror.NewError(err, 403, "invalid token", map[string]interface{}{"opration": "VerifyToken"})
}
