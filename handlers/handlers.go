package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func HandleRoot(c *fiber.Ctx) error {
	return c.JSON(struct {
		Msg string `json:"Msg"`
	}{Msg: "salam jahaaaaaaaaaaaaaaaaan!"})
}
