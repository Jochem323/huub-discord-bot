package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Connect to Postgres database
func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable host=" + os.Getenv("POSTGRES_HOST") + " port=" + os.Getenv("POSTGRES_PORT")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to Postgres!")

	return &PostgresStore{
		db: db,
	}, nil
}

// Initialize the database
func (s *PostgresStore) Init() error {
	err := s.CreatekeywordTable()
	if err != nil {
		return err
	}

	err = s.CreateguildTable()
	if err != nil {
		return err
	}

	return nil
}
