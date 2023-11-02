package handlers

import (
	"encoding/json"
	"fmt"
	"hangman-api/database"
	"hangman-api/models"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func getRandomWord() (string, error) {
	resp, err := http.Get("https://random-word-api.herokuapp.com/word?lang=en")
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var words []string
	json.Unmarshal(body, &words)
	return words[0], nil
}

func HandleCreateGame(c *fiber.Ctx) error {

	userIdString := fmt.Sprintf("%v", c.Locals("userId"))

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		return WriteResponse(c, http.StatusUnauthorized, "forbidden resources")
	}

	secret_word, err := getRandomWord()
	thisUser, err := database.GetUserByID(userId)

	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}

	user := models.ConvertToUserDto(thisUser)
	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}
	thisGame, err := database.CreateGame(secret_word, userId)
	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}

	game, err := models.ConvertGameIntoGamesDto(thisGame, user)
	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}

	return c.JSON(game)
}

func HandleGetSingleGame(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return WriteResponse(c, http.StatusNotFound, "game not found")
	}

	thisGame, err := database.GetSingleGameById(id)
	if err != nil {
		return WriteResponse(c, http.StatusNotFound, "game not found")
	}
	thisUser, err := database.GetUserByID(thisGame.User_Id)
	if err != nil {
		return WriteResponse(c, http.StatusNotFound, "user not found")
	}
	user := models.ConvertToUserDto(thisUser)
	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}

	game, err := models.ConvertGameIntoGamesDto(thisGame, user)
	if err != nil {
		return WriteResponse(c, http.StatusUnauthorized, "forbidden resources")
	}

	return c.JSON(game)

}

func HandleGetMyGames(c *fiber.Ctx) error {

	userIdString := fmt.Sprintf("%v", c.Locals("userId"))

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		return WriteResponse(c, http.StatusUnauthorized, "forbidden resources")
	}
	body := new(models.GetMyGamesDto)
	if err := c.BodyParser(body); err != nil {
		body.Page = 0
		body.Limit = 10
	}

	thisUser, err := database.GetUserByID(userId)
	if err != nil {
		return WriteResponse(c, http.StatusUnauthorized, "forbidden resources")
	}
	user := models.ConvertToUserDto(thisUser)
	MyGames, err := database.GetSingleUsersGames(userId, body.Page, body.Limit)
	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}

	var games []models.GameDto

	for _, s := range MyGames {
		game, err := models.ConvertGameIntoGamesDto(s, user)
		if err != nil {
			return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
		}

		games = append(games, *game)
	}

	return c.Status(http.StatusOK).JSON(games)
}

func HandleGuessLetter(c *fiber.Ctx) error {

	userIdString := fmt.Sprintf("%v", c.Locals("userId"))

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		return WriteResponse(c, http.StatusUnauthorized, "forbidden resources")
	}
	body := new(models.GuessLetterDto)
	if err := c.BodyParser(body); err != nil {
		return WriteResponse(c, http.StatusBadRequest, "please provide valid information")
	}

	if len(body.Char) > 1 {
		return WriteResponse(c, http.StatusBadRequest, "please provide only a single letter")
	}

	thisUser, err := database.GetUserByID(userId)
	if err != nil {
		return WriteResponse(c, http.StatusUnauthorized, "forbidden resources")
	}

	thisGame, err := database.GetSingleGameById(body.GameId)

	if err != nil {
		return WriteResponse(c, http.StatusNotFound, "bad request: no such game found")
	}

	if thisGame.Chances < 1 || thisGame.Status != "ongoing" {
		return WriteResponse(c, http.StatusBadRequest, "the game is over, please start a new game")
	}

	if exists := strings.Contains(thisGame.Guessed_letters, body.Char); exists {
		return WriteResponse(c, http.StatusBadRequest, "you have guessed this letter before, try a new character")
	}

	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}
	if thisGame.User_Id != thisUser.ID {
		return WriteResponse(c, fiber.StatusUnauthorized, "unathorized")
	}

	err2 := database.HandleGuessLetter(thisGame, body.Char)
	if err2 != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}

	user := models.ConvertToUserDto(thisUser)
	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}

	game, err := models.ConvertGameIntoGamesDto(thisGame, user)
	if err != nil {
		return WriteResponse(c, http.StatusInternalServerError, "this one is on us")
	}

	if game.Chances == 0 {
		return WriteResponse(c, http.StatusOK, fmt.Sprintf("you lost, the correct word was : %v", thisGame.Secret_word))
	}
	if game.Status == "won" {
		return WriteResponse(c, http.StatusOK, "congatulations, you won. start a new game now")
	}

	return c.Status(http.StatusOK).JSON(game)
}

func HandleGameStatistics(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")

	if err != nil {
		return fmt.Errorf("bad request")
	}
	r := database.FetchStatistics(userId)

	return c.Status(http.StatusOK).JSON(r)
}
