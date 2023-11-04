package models

import "time"

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
