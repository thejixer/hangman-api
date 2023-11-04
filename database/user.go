package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/thewhiterabbit1994/hangman-api/models"
	"github.com/thewhiterabbit1994/hangman-api/utils"
)

func (s *PostgresStore) CreateUserTable() error {

	query := `create table if not exists users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100),
		password VARCHAR,
		created_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateUser(Name, Email, Password string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(Password)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		Name:      Name,
		Email:     Email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}

	query := `
	INSERT INTO USERS (name, email, password, created_at)
	VALUES ($1, $2, $3, $4) RETURNING id`
	lastInsertId := 0

	insertErr := s.db.QueryRow(query, newUser.Name, newUser.Email, newUser.Password, newUser.CreatedAt).Scan(&lastInsertId)
	if insertErr != nil {
		return nil, insertErr
	}
	newUser.ID = lastInsertId

	return newUser, nil
}

func (s *PostgresStore) GetUserByEmail(email string) (*models.User, error) {
	rows, err := s.db.Query("SELECT * FROM USERS WHERE EMAIL = $1", email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoUsers(rows)
	}

	return nil, errors.New("not found")
}

func (s *PostgresStore) GetUserByID(id int) (*models.User, error) {
	rows, err := s.db.Query("SELECT * FROM USERS WHERE ID = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoUsers(rows)
	}

	return nil, errors.New("not found")
}

func (s *PostgresStore) GetUsers(page, limit int) ([]*models.User, error) {
	offset := page * limit
	query := "SELECT * FROM USERS ORDER BY id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY"
	rows, err := s.db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}
	users := []*models.User{}
	for rows.Next() {
		u, err := scanIntoUsers(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil

}

func scanIntoUsers(rows *sql.Rows) (*models.User, error) {
	u := new(models.User)
	if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return nil, err
	}
	return u, nil
}
