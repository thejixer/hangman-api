package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	dataprocesslayer "github.com/thewhiterabbit1994/hangman-api/data-process-layer"
	"github.com/thewhiterabbit1994/hangman-api/models"
)

func (h *HandlerService) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	body := new(models.PaginationDto)

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		body.Limit = 10
	}

	users, err := h.db.GetUsers(body.Page, body.Limit)

	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, "this one is on us")
		return
	}

	var result []models.UserDto

	for _, s := range users {
		result = append(result, dataprocesslayer.ConvertToUserDto(s))
	}

	WriteJSON(w, http.StatusOK, result)

}

func (h *HandlerService) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userIdString := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		WriteResponse(w, http.StatusNotFound, "user not found")
		return
	}

	thisUser, err := h.db.GetUserByID(userId)
	if err != nil {
		WriteResponse(w, http.StatusNotFound, "user not found")
		return
	}

	user := dataprocesslayer.ConvertToUserDto(thisUser)

	WriteJSON(w, http.StatusOK, user)
}

func (h *HandlerService) HandleGetUserStatistics(w http.ResponseWriter, r *http.Request) {

	userIdString := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		WriteResponse(w, http.StatusNotFound, "user not found")
		return
	}
	statistics := h.db.FetchStatistics(userId)

	WriteJSON(w, http.StatusOK, statistics)
}
