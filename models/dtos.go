package models

type Response struct {
	Msg        string `json:"msg"`
	StatusCode int    `json:"statusCode"`
}

type TokenDto struct {
	Token string `json:"token"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUsersDto struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type GetMyGamesDto struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type GuessLetterDto struct {
	GameId int    `json:"gameId"`
	Char   string `json:"char"`
}

type StatisticsDto struct {
	WonCount     int `json:"wonCount"`
	LostCount    int `json:"lostCount"`
	OngoingCount int `json:"ongoingCount"`
}
