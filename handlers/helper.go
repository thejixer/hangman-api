package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/thewhiterabbit1994/hangman-api/models"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func WriteResponse(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, models.ResponseDto{Msg: msg, StatusCode: status})
}

func GetUserIdByToken(r *http.Request) (int, error) {
	userIdString := r.Context().Value("userId").(string)
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
