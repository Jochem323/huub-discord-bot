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
		comment TEXT,
		created_by TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		active BOOLEAN NOT NULL,
		revoked BOOLEAN NOT NULL,
		ratelimit BOOLEAN NOT NULL
	);`

	_, err := s.db.Exec(query)
	return err
}

func ScanIntoAPIKey(rows *sql.Rows) (any, error) {
	key := new(common.APIKey)
	err := rows.Scan(
		&key.ID,
		&key.Admin,
		&key.GuildID,
		&key.Comment,
		&key.CreatedBy,
		&key.CreatedAt,
		&key.Active,
		&key.Revoked,
		&key.Ratelimit,
	)

	return *key, err
}

func (s *PostgresStore) GetKeys() ([]common.APIKey, error) {
	query := `SELECT * FROM api_keys;`

	results, err := s.GetMultiple(ScanIntoAPIKey, query)
	if err != nil {
		return nil, err
	}

	keys := []common.APIKey{}
	for _, result := range results {
		key := result.(common.APIKey)
		keys = append(keys, key)
	}

	return keys, nil
}

func (s *PostgresStore) GetKey(keyID int) (common.APIKey, error) {
	query := `SELECT * FROM api_keys WHERE id = $1;`

	result, err := s.GetOne(ScanIntoAPIKey, query, keyID)
	if err != nil {
		return common.APIKey{}, err
	}

	key := result.(common.APIKey)

	return key, nil
}

func (s *PostgresStore) AddKey(key common.APIKey) (int, error) {
	query := `INSERT INTO api_keys
	(admin, guild_id, comment, created_by, created_at, active, revoked, ratelimit)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	;`

	id, err := s.ExecReturnId(query, key.Admin, key.GuildID, key.Comment, key.CreatedBy, key.CreatedAt, key.Active, key.Revoked, key.Ratelimit)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *PostgresStore) UpdateKey(key common.APIKey) error {
	query := `UPDATE api_keys SET comment = $1, active = $2, revoked = $3, ratelimit = $4 WHERE id = $5;`
	return s.Exec(query, key.Comment, key.Active, key.Revoked, key.Ratelimit, key.ID)
}

func (s *PostgresStore) DeleteKey(keyID int) error {
	query := `DELETE FROM api_keys WHERE id = $1;`
	return s.Exec(query, keyID)
}
