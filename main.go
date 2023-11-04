package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thewhiterabbit1994/hangman-api/database"
	"github.com/thewhiterabbit1994/hangman-api/handlers"
	"github.com/thewhiterabbit1994/hangman-api/server"
)

func init() {
	godotenv.Load()
}

func main() {

	dbName := os.Getenv("DB_NAME")

	store, err := database.NewPostgresStore(dbName)

	if err != nil {
		log.Fatal("couldnt reach the database :", err)
	}
	store.Init()

	handlerService := handlers.NewHandlerService(store)

	if err != nil {
		log.Fatal("db a'int reachable")
	}

	server := server.NewAPIServer(":3000", handlerService)
	server.Run()
}
