package database

import (
	"database/sql"
	"hangman-api/models"

	_ "github.com/lib/pq"
)

var store *PostgresStore

type Storage interface {
	CreateUser(*models.User) error
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	s.createGameTable()
	return s.CreateUserTable()
}

func NewPostgresStore() (*PostgresStore, error) {
	conString := "user=postgres dbname=hangman password=postgres sslmode=disable"
	db, err := sql.Open("postgres", conString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store = &PostgresStore{
		db: db,
	}

	return &PostgresStore{
		db: db,
	}, nil
}
