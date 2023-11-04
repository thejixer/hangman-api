package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/thewhiterabbit1994/hangman-api/models"
)

func (s *PostgresStore) createGameTable() error {

	query := `create table if not exists games (
		id SERIAL PRIMARY KEY,
		userId INT,
		sercet_word VARCHAR(100),
		guessedLetters VARCHAR(100),
		chances INT DEFAULT 5,
		status VARCHAR(50),
		created_at TIMESTAMP,
		finished_at TIMESTAMP
		)`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	query2 := `ALTER TABLE games ADD FORIEGN KEY ("userId") REFRENCES "USERS ("id")"`
	_, err2 := s.db.Exec(query2)

	return err2
}

func (s *PostgresStore) CreateGame(Secret_word string, User_Id int) (*models.Game, error) {

	thisGame := models.Game{
		User_Id:         User_Id,
		Secret_word:     Secret_word,
		Guessed_letters: "",
		Chances:         5,
		Status:          "ongoing",
		Created_at:      time.Now().UTC(),
		Finished_at:     time.Time{},
	}

	query := `
	INSERT INTO GAMES (userId, sercet_word, guessedLetters, status, created_at, finished_at)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	lastInsertId := 0
	insertErr := s.db.QueryRow(query, thisGame.User_Id, thisGame.Secret_word, thisGame.Guessed_letters, thisGame.Status, thisGame.Created_at, thisGame.Finished_at).Scan(&lastInsertId)

	if insertErr != nil {
		return nil, insertErr
	}
	thisGame.ID = lastInsertId

	return &thisGame, nil

}

func (s *PostgresStore) GetSingleGameById(id int) (*models.Game, error) {
	rows, err := s.db.Query("SELECT * FROM GAMES WHERE ID = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoGames(rows)
	}

	return nil, errors.New("not found")
}

func (s *PostgresStore) GetSingleUsersGames(id, page, limit int) ([]*models.Game, error) {
	offset := page * limit
	query := "SELECT * FROM GAMES WHERE userId = $1 ORDER BY id OFFSET $2 ROWS FETCH NEXT $3 ROWS ONLY"
	rows, err := s.db.Query(query, id, offset, limit)
	if err != nil {
		return nil, err
	}
	games := []*models.Game{}
	for rows.Next() {
		g, err := scanIntoGames(rows)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, nil
}

func (s *PostgresStore) HandleGuessLetter(g *models.Game, character string) error {

	query := `
		UPDATE GAMES
		SET guessedLetters = $1, chances = $2, status = $3
		WHERE id = $4;
	`

	_, err := s.db.Exec(query, g.Guessed_letters, g.Chances, g.Status, g.ID)

	if err != nil {
		return err
	}

	return nil

}

func getStatistics(s *PostgresStore, str string, id int, target *int, wg *sync.WaitGroup) {

	defer wg.Done()

	query := `SELECT * FROM GAMES WHERE status = $1 AND userId = $2`

	rows, err := s.db.Query(query, str, id)
	if err != nil {
		fmt.Println("err : ", err)
	}

	var count int
	for rows.Next() {
		count++
	}

	*target = count
}

func (s *PostgresStore) FetchStatistics(id int) models.StatisticsDto {

	var wg sync.WaitGroup
	wg.Add(3)

	r := models.StatisticsDto{}
	go getStatistics(s, "won", id, &r.WonCount, &wg)
	go getStatistics(s, "lost", id, &r.LostCount, &wg)
	go getStatistics(s, "ongoing", id, &r.OngoingCount, &wg)

	wg.Wait()

	return r

}

func scanIntoGames(rows *sql.Rows) (*models.Game, error) {
	g := new(models.Game)
	if err := rows.Scan(&g.ID, &g.User_Id, &g.Secret_word, &g.Guessed_letters, &g.Chances, &g.Status, &g.Created_at, &g.Finished_at); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return g, nil
}
