package storage

import (
	"database/sql"
	"huub-discord-bot/common"
)

func (s *PostgresStore) CreateAPIKeyTable() error {
	query := `CREATE TABLE IF NOT EXISTS api_keys (
		id SERIAL PRIMARY KEY,
		admin BOOLEAN NOT NULL,
		guild_id TEXT,
		created_by TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		active BOOLEAN NOT NULL,
		revoked BOOLEAN NOT NULL,
		ratelimit BOOLEAN NOT NULL
	);`

	_, err := s.db.Exec(query)
	return err
}

func ScanRowsIntoAPIKey(rows *sql.Rows) (*common.APIKey, error) {
	key := new(common.APIKey)
	err := rows.Scan(
		&key.ID,
		&key.Admin,
		&key.GuildID,
		&key.CreatedBy,
		&key.CreatedAt,
		&key.Active,
		&key.Revoked,
		&key.Ratelimit,
	)

	return key, err
}

func ScanRowIntoAPIKey(row *sql.Row) (*common.APIKey, error) {
	key := new(common.APIKey)
	err := row.Scan(
		&key.ID,
		&key.Admin,
		&key.GuildID,
		&key.CreatedBy,
		&key.CreatedAt,
		&key.Active,
		&key.Revoked,
		&key.Ratelimit,
	)

	return key, err
}

func (s *PostgresStore) GetKeys() (*[]common.APIKey, error) {
	query := `SELECT * FROM api_keys;`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keys := []common.APIKey{}
	for rows.Next() {
		key, err := ScanRowsIntoAPIKey(rows)
		if err != nil {
			return nil, err
		}

		keys = append(keys, *key)
	}

	return &keys, nil
}

func (s *PostgresStore) GetKey(keyID int) (*common.APIKey, error) {
	query := `SELECT * FROM api_keys WHERE id = $1;`

	row := s.db.QueryRow(query, keyID)

	key, err := ScanRowIntoAPIKey(row)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s *PostgresStore) AddKey(key *common.APIKey) (int, error) {
	query := `INSERT INTO api_keys (admin, guild_id, created_by, created_at, active, revoked, ratelimit) VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := s.db.Exec(query, key.Admin, key.GuildID, key.CreatedBy, key.CreatedAt, key.Active, key.Revoked, key.Ratelimit)

	if err != nil {
		return -1, err
	}

	query = `SELECT id FROM api_keys ORDER BY id DESC LIMIT 1;`

	row := s.db.QueryRow(query)
	if err != nil {
		return -1, err
	}

	var id int
	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *PostgresStore) UpdateKey(key *common.APIKey) error {
	query := `UPDATE api_keys SET active = $1, revoked = $2, ratelimit = $3 WHERE id = $4;`

	_, err := s.db.Exec(query, key.Active, key.Revoked, key.Ratelimit, key.ID)
	return err
}

func (s *PostgresStore) DeleteKey(keyID int) error {
	query := `DELETE FROM api_keys WHERE id = $1;`

	_, err := s.db.Exec(query, keyID)
	return err
}
