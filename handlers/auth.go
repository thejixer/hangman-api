package handlers

import (
	"fmt"
	"hangman-api/database"
	"hangman-api/models"
	"hangman-api/utils"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func HandleSignUp(c *fiber.Ctx) error {
	body := new(models.SignUpDto)

	if err := c.BodyParser(body); err != nil {
		return err
	}

	if len(body.Name) < 3 || len(body.Email) < 3 || len(body.Password) < 3 {
		return WriteResponse(c, http.StatusBadRequest, "bad request: insufficent data")
	}

	user, err := database.NewUser(body.Name, body.Email, body.Password)

	if err != nil {
		return err
	}

	tokenString, err := utils.SignToken(user.ID)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(models.TokenDto{Token: tokenString})
}

func HandleMe(c *fiber.Ctx) error {
	userIdString := fmt.Sprintf("%v", c.Locals("userId"))

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		return WriteResponse(c, http.StatusUnauthorized, "forbidden resources")
	}

	thisUser, err := database.GetUserByID(userId)

	result := models.ConvertToUserDto(thisUser)

	return c.JSON(result)
}

func HandleLogin(c *fiber.Ctx) error {
	body := new(models.SignUpDto)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	if len(body.Email) < 3 || len(body.Password) < 3 {
		return WriteResponse(c, http.StatusBadRequest, "bad request: insufficent data")
	}

	thisUser, err := database.GetUserByEmail(body.Email)
	if err != nil {
		fmt.Println(err)
		return WriteResponse(c, http.StatusBadRequest, "bad request: no such user found")
	}

	if match := utils.CheckPasswordHash(body.Password, thisUser.Password); !match {
		return WriteResponse(c, http.StatusBadRequest, "bad request: password doesnt match")
	}

	tokenString, err := utils.SignToken(thisUser.ID)

	return c.JSON(models.TokenDto{Token: tokenString})

}
