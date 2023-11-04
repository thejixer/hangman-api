package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	dataprocesslayer "github.com/thewhiterabbit1994/hangman-api/data-process-layer"
	"github.com/thewhiterabbit1994/hangman-api/models"
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

func (h *HandlerService) HandleCreateGame(w http.ResponseWriter, r *http.Request) {

	userId, err := GetUserIdByToken(r)
	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, "forbidden resources")
		return
	}

	thisUser, err := h.db.GetUserByID(userId)
	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, "forbidden resources")
		return
	}

	secret_word, err := getRandomWord()
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, "this one is on us: error fetching random word from external api")
		return
	}

	user := dataprocesslayer.ConvertToUserDto(thisUser)
	thisGame, err := h.db.CreateGame(secret_word, userId)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, "this one is on us")
		return
	}

	game := dataprocesslayer.ConvertGameIntoGamesDto(thisGame, user)

	WriteJSON(w, http.StatusOK, game)
}

func (h *HandlerService) HandleGetSingleGame(w http.ResponseWriter, r *http.Request) {
	gameIdString := chi.URLParam(r, "id")
	gameId, err := strconv.Atoi(gameIdString)

	if err != nil {
		WriteResponse(w, http.StatusNotFound, "game not found")
		return
	}

	thisGame, err := h.db.GetSingleGameById(gameId)
	if err != nil {
		WriteResponse(w, http.StatusNotFound, "game not found")
		return
	}
	thisUser, err := h.db.GetUserByID(thisGame.User_Id)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, "this one is on us")
		return
	}

	user := dataprocesslayer.ConvertToUserDto(thisUser)
	game := dataprocesslayer.ConvertGameIntoGamesDto(thisGame, user)

	WriteJSON(w, http.StatusOK, game)
}

func (h *HandlerService) HandleGetMyGames(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserIdByToken(r)
	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, "forbidden resources")
		return
	}
	body := new(models.PaginationDto)

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		body.Limit = 10
	}

	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, "forbidden resources")
		return
	}

	thisUser, err := h.db.GetUserByID(userId)
	user := dataprocesslayer.ConvertToUserDto(thisUser)

	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, "forbidden resources")
		return
	}

	MyGames, err := h.db.GetSingleUsersGames(userId, body.Page, body.Limit)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, "this one is on us")
		return
	}

	var games = []models.GameDto{}

	for _, s := range MyGames {
		game := dataprocesslayer.ConvertGameIntoGamesDto(s, user)

		games = append(games, *game)
	}

	WriteJSON(w, http.StatusOK, games)
}

func (h *HandlerService) HandleGuessLetter(w http.ResponseWriter, r *http.Request) {

	fmt.Println("yes")
	userId, err := GetUserIdByToken(r)
	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, "forbidden resources")
		return
	}

	body := new(models.GuessLetterDto)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		WriteResponse(w, http.StatusBadRequest, "bad request: insufficient data")
		return
	}
	if len(body.Char) > 1 {
		WriteResponse(w, http.StatusBadRequest, "please provide only a single letter")
		return
	}
	thisUser, err := h.db.GetUserByID(userId)
	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, "forbidden resources")
		return
	}
	thisGame, err := h.db.GetSingleGameById(body.GameId)
	if err != nil {
		WriteResponse(w, http.StatusNotFound, "bad request: no such game found")
		return
	}
	if thisGame.User_Id != thisUser.ID {
		WriteResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}

	if thisGame.Chances < 1 || thisGame.Status != "ongoing" {
		WriteResponse(w, http.StatusBadRequest, "the game is over, please start a new game")
		return
	}

	if exists := strings.Contains(thisGame.Guessed_letters, body.Char); exists {
		WriteResponse(w, http.StatusBadRequest, "you have guessed this letter before, try a new character")
		return
	}

	exists := strings.Contains(thisGame.Secret_word, body.Char)

	thisGame.Guessed_letters += body.Char

	if !exists {
		thisGame.Chances--
	}

	if thisGame.Chances < 1 {
		thisGame.Status = "lost"
	}

	guessLetters := dataprocesslayer.BuildGuessedLetterJSON(thisGame)
	dashifiedString := dataprocesslayer.DashifyString(thisGame.Secret_word, guessLetters)

	if hasDash := strings.Contains(dashifiedString, "-"); !hasDash {
		thisGame.Status = "won"
	}

	updateErr := h.db.HandleGuessLetter(thisGame, body.Char)

	if updateErr != nil {
		WriteResponse(w, http.StatusInternalServerError, "this one's on us")
		return
	}

	user := dataprocesslayer.ConvertToUserDto(thisUser)
	game := dataprocesslayer.ConvertGameIntoGamesDto(thisGame, user)

	if game.Chances == 0 {
		WriteResponse(w, http.StatusOK, fmt.Sprintf("you lost, the correct word was : %v", thisGame.Secret_word))
		return
	}
	if game.Status == "won" {
		WriteResponse(w, http.StatusOK, "congatulations, you won. start a new game now")
		return
	}

	WriteJSON(w, http.StatusOK, game)

}
