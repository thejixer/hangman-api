package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

func SignToken(id int) (string, error) {

	secret := os.Getenv("JWT_SECRET")
	mySigningKey := []byte(secret)

	// Create the Claims
	claims := jwt.MapClaims{
		"id": fmt.Sprintf("%v", id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)

}
