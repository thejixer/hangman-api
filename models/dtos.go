package models

type ResponseDto struct {
	Msg        string `json:"msg"`
	StatusCode int    `json:"statusCode"`
}

type SignUpDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenDto struct {
	Token string `json:"token"`
}

type PaginationDto struct {
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
