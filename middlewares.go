package main

import (
	"fmt"
	"hangman-api/models"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JwtAuth(c *fiber.Ctx) error {
	token := c.Get("auth")
	c.Locals("token", token)
	return c.Next()
}

func AuthGaurd(c *fiber.Ctx) error {
	tokenString := fmt.Sprintf("%v", c.Locals("token"))
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusForbidden).JSON(models.Response{Msg: "forbidden resources", StatusCode: http.StatusForbidden})
	}

	if !token.Valid {
		return c.Status(http.StatusForbidden).JSON(models.Response{Msg: "forbidden resources", StatusCode: http.StatusForbidden})
	}

	claims := token.Claims.(jwt.MapClaims)

	c.Locals("userId", claims["id"])

	return c.Next()
}
