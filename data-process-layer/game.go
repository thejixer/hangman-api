package dataprocesslayer

import (
	"strings"

	"github.com/thewhiterabbit1994/hangman-api/models"
	"github.com/thewhiterabbit1994/hangman-api/utils"
)

func ConvertGameIntoGamesDto(g *models.Game, u models.UserDto) *models.GameDto {

	guessedLetters := BuildGuessedLetterJSON(g)
	dashifiedString := DashifyString(g.Secret_word, guessedLetters)

	return &models.GameDto{
		ID:              g.ID,
		User:            u,
		Secret_word:     dashifiedString,
		Guessed_letters: guessedLetters,
		Chances:         g.Chances,
		Status:          g.Status,
		Created_at:      g.Created_at,
		Finished_at:     g.Finished_at,
	}
}

func DashifyString(s string, guessLetters []models.GuessedLetter) string {

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

func BuildGuessedLetterJSON(g *models.Game) []models.GuessedLetter {

	guessLetters := []models.GuessedLetter{}

	for _, letter := range g.Guessed_letters {

		character := string(letter)
		var x models.GuessedLetter
		x.Char = character

		if exists := strings.Contains(g.Secret_word, character); exists {
			x.Contains = true
		}

		guessLetters = append(guessLetters, x)
	}

	return guessLetters
}
