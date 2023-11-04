package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thewhiterabbit1994/hangman-api/handlers"
)

type APIServer struct {
	listenAddr     string
	handlerService *handlers.HandlerService
}

func NewAPIServer(listenAddr string, handlerService *handlers.HandlerService) *APIServer {

	return &APIServer{
		listenAddr:     listenAddr,
		handlerService: handlerService,
	}
}

func (s *APIServer) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	s.applyRoutes(r)
	fmt.Println("app is running on port ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, r)
}

func (s *APIServer) applyRoutes(r *chi.Mux) {

	// protected routes
	r.Group(func(r chi.Router) {

		r.Use(s.handlerService.AuthGaurd)

		r.Post("/auth/me", s.handlerService.HandleMe)
		r.Post("/game/", s.handlerService.HandleCreateGame)
		r.Post("/game/mygames", s.handlerService.HandleGetMyGames)
		r.Post("/game/guess", s.handlerService.HandleGuessLetter)
	})

	// public routes
	r.Get("/", s.handlerService.HandleHelloWorld)
	r.Post("/auth/signup", s.handlerService.HandleSignUp)
	r.Post("/auth/login", s.handlerService.HandleLogin)
	r.Post("/user/", s.handlerService.HandleGetUsers)
	r.Post("/user/{id}", s.handlerService.HandleGetUser)
	r.Get("/game/{id}", s.handlerService.HandleGetSingleGame)
	r.Post("/user/statistics/{id}", s.handlerService.HandleGetUserStatistics)

}
