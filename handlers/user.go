package handlers

import (
	"hangman-api/database"
	"hangman-api/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {

	body := new(models.GetUsersDto)
	if err := c.BodyParser(body); err != nil {
		body.Page = 0
		body.Limit = 10
	}

	users, err := database.GetUsers(body.Page, body.Limit)

	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}

	var result []models.UserDto

	for _, s := range users {
		result = append(result, models.ConvertToUserDto(s))
	}

	return c.Status(http.StatusOK).JSON(result)
}

func HandleSingleUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return WriteResponse(c, http.StatusNotFound, "user not found")
	}

	thisUser, err := database.GetUserByID(id)
	if err != nil {
		return WriteResponse(c, http.StatusNotFound, "user not found")

	}
	user := models.ConvertToUserDto(thisUser)
	return c.Status(http.StatusOK).JSON(user)
}
