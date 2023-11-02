package handlers

import (
	"hangman-api/models"

	"github.com/gofiber/fiber/v2"
)

func WriteResponse(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(models.Response{Msg: msg, StatusCode: status})
}
