package handlers

import (
	"net/http"

	"github.com/thewhiterabbit1994/hangman-api/database"
)

type HandlerService struct {
	db *database.PostgresStore
}

func NewHandlerService(store *database.PostgresStore) *HandlerService {
	return &HandlerService{
		db: store,
	}
}

func (h *HandlerService) HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
