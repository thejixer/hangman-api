package main

import (
	"flag"
	"hangman-api/database"
	"hangman-api/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {

	// Parse command-line flags
	flag.Parse()
	store, err := database.NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	store.Init()

	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})
	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	app.Use(JwtAuth)
	auth := app.Group("/auth")
	user := app.Group("/user")
	game := app.Group("/game")

	auth.Post("/signup", handlers.HandleSignUp)
	auth.Post("/login", handlers.HandleLogin)
	auth.Post("/me", AuthGaurd, handlers.HandleMe)

	user.Post("/", handlers.HandleGetUsers)
	user.Post("/:id", handlers.HandleSingleUser)
	user.Post("/statistics/:id", AuthGaurd, handlers.HandleGameStatistics)

	game.Post("/", AuthGaurd, handlers.HandleCreateGame)
	game.Post("/mygames", AuthGaurd, handlers.HandleGetMyGames)
	game.Get("/:id", handlers.HandleGetSingleGame)
	game.Post("/guess", AuthGaurd, handlers.HandleGuessLetter)

	log.Fatal(app.Listen(*port)) // go run app.go -port=:3000

}
