package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/thewhiterabbit1994/hangman-api/models"
)

type Storage interface {
	CreateUser(Name, Email, Password string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUsers(page, limit int) ([]*models.User, error)
	CreateGame(Secret_word string, User_Id int) (*models.Game, error)
	GetSingleGameById(id int) (*models.Game, error)
	GetSingleUsersGames(id, page, limit int) ([]*models.Game, error)
	HandleGuessLetter(g *models.Game, character string) error
	FetchStatistics(id int) models.StatisticsDto
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(dbName string) (*PostgresStore, error) {

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	conString := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=disable", dbUser, dbName, dbPassword)
	db, err := sql.Open("postgres", conString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	s.createGameTable()
	return s.CreateUserTable()
}
