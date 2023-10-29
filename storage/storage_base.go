package storage

import (
	"database/sql"
	"log"
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

	return &PostgresStore{
		db: db,
	}, nil
}

// Initialize the database
func (s *PostgresStore) Init() error {
	s.log = log.New(os.Stdout, "Storage: ", log.Ldate|log.Ltime)

	err := s.CreateKeywordTable()
	if err != nil {
		return err
	}

	err = s.CreateGuildTable()
	if err != nil {
		return err
	}

	err = s.CreateAPIKeyTable()
	if err != nil {
		return err
	}

	s.log.Println("Database initialized")

	return nil
}

func (s *PostgresStore) GetMultiple(scanFunc func(*sql.Rows) (any, error), query string, args ...any) ([]any, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []any
	for rows.Next() {
		result, err := scanFunc(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *PostgresStore) GetOne(scanFunc func(*sql.Rows) (any, error), query string, args ...any) (any, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, ErrNotFound
	}

	result, err := scanFunc(rows)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return nil, ErrMultipleFound
	}

	return result, nil
}

func (s *PostgresStore) Exec(query string, args ...any) error {
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *PostgresStore) ExecReturnId(query string, args ...any) (int, error) {
	var id int
	err := s.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
