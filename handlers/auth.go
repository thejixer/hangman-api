package handlers

import (
	"encoding/json"
	"net/http"

	dataprocesslayer "github.com/thewhiterabbit1994/hangman-api/data-process-layer"
	"github.com/thewhiterabbit1994/hangman-api/models"
	"github.com/thewhiterabbit1994/hangman-api/utils"
)

func (h *HandlerService) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	body := new(models.SignUpDto)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		WriteResponse(w, http.StatusBadRequest, "bad request: please provide sufficient data")
		return
	}

	if len(body.Name) < 3 || len(body.Email) < 3 || len(body.Password) < 3 {
		WriteResponse(w, http.StatusBadRequest, "bad request: insufficent data")
		return
	}

	thisUser, _ := h.db.GetUserByEmail(body.Email)

	if thisUser != nil {
		WriteResponse(w, http.StatusBadRequest, "bad request: this email already exists in the database")
		return
	}

	user, err := h.db.CreateUser(body.Name, body.Email, body.Password)

	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	tokenString, err := utils.SignToken(user.ID)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, models.TokenDto{Token: tokenString})
}

func (h *HandlerService) HandleMe(w http.ResponseWriter, r *http.Request) {
	userId, err := GetUserIdByToken(r)
	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, "forbidden resources")
		return
	}

	thisUser, err := h.db.GetUserByID(userId)

	user := dataprocesslayer.ConvertToUserDto(thisUser)
	WriteJSON(w, http.StatusOK, user)

}

func (h *HandlerService) HandleLogin(w http.ResponseWriter, r *http.Request) {
	body := new(models.LoginDto)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		WriteResponse(w, http.StatusBadRequest, "bad request: please provide sufficient data")
		return
	}

	if len(body.Email) < 3 || len(body.Password) < 3 {
		WriteResponse(w, http.StatusBadRequest, "bad request: insufficent data")
		return
	}

	thisUser, err := h.db.GetUserByEmail(body.Email)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, "bad request: no such user found")
		return
	}

	if match := utils.CheckPasswordHash(body.Password, thisUser.Password); !match {
		WriteResponse(w, http.StatusBadRequest, "bad request: password doesnt match")
		return
	}

	tokenString, err := utils.SignToken(thisUser.ID)

	WriteJSON(w, http.StatusOK, models.TokenDto{Token: tokenString})

}
