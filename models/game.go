package models

import (
	"hangman-api/utils"
	"strings"
	"time"
)

type GuessedLetter struct {
	Char     string `json:"char"`
	Contains bool   `json:"contains"`
}

type Game struct {
	ID              int       `json:"id"`
	User_Id         int       `json:"userId"`
	Secret_word     string    `json:"secret_word"`
	Guessed_letters string    `json:"guessed_letters"`
	Chances         int       `json:"chances"`
	Status          string    `json:"status"`
	Created_at      time.Time `json:"created_at"`
	Finished_at     time.Time `json:"finished_at"`
}

type GameDto struct {
	ID              int             `json:"id"`
	User            UserDto         `json:"user"`
	Secret_word     string          `json:"secret_word"`
	Guessed_letters []GuessedLetter `json:"guessed_letters"`
	Chances         int             `json:"chances"`
	Status          string          `json:"status"`
	Created_at      time.Time       `json:"created_at"`
	Finished_at     time.Time       `json:"finished_at"`
}

func DashifyString(s string, guessLetters []GuessedLetter) string {

	var trueLetters []string

	for _, char := range guessLetters {
		if char.Contains {
			trueLetters = append(trueLetters, char.Char)
		}
	}

	str := ""

	for _, char := range s {
		c := string(char)
		if x := utils.Contains(trueLetters, c); x {
			str += c
		} else {
			str += "-"
		}
	}
	return str
}

func BuildGuessedLetterJSON(g *Game) []GuessedLetter {

	guessLetters := []GuessedLetter{}

	for _, letter := range g.Guessed_letters {

		character := string(letter)
		var x GuessedLetter
		x.Char = character

		if exists := strings.Contains(g.Secret_word, character); exists {
			x.Contains = true
		}

		guessLetters = append(guessLetters, x)
	}

	return guessLetters
}

func ConvertGameIntoGamesDto(g *Game, u UserDto) (*GameDto, error) {

	guessedLetters := BuildGuessedLetterJSON(g)
	dashifiedString := DashifyString(g.Secret_word, guessedLetters)

	return &GameDto{
		ID:              g.ID,
		User:            u,
		Secret_word:     dashifiedString,
		Guessed_letters: guessedLetters,
		Chances:         g.Chances,
		Status:          g.Status,
		Created_at:      g.Created_at,
		Finished_at:     g.Finished_at,
	}, nil
}
