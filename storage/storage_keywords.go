package storage

import (
	"database/sql"
	"huub-discord-bot/common"
)

func (s *PostgresStore) CreateKeywordTable() error {
	query := `CREATE TABLE IF NOT EXISTS keywords (
		id SERIAL PRIMARY KEY,
		guild_id TEXT NOT NULL,
		keyword TEXT NOT NULL,
		reaction TEXT NOT NULL
	);`

	_, err := s.db.Exec(query)
	return err
}

func ScanIntoKeyword(rows *sql.Rows) (*common.Keyword, error) {
	keyword := new(common.Keyword)
	err := rows.Scan(
		&keyword.ID,
		&keyword.GuildID,
		&keyword.Keyword,
		&keyword.Reaction,
	)

	return keyword, err
}

func (s *PostgresStore) GetKeywords(guildID string) (*[]common.Keyword, error) {
	query := `SELECT * FROM keywords WHERE guild_id = $1;`

	rows, err := s.db.Query(query, guildID)
	if err != nil {
		return nil, err
	}

	keywords := []common.Keyword{}
	for rows.Next() {
		keyword, err := ScanIntoKeyword(rows)
		if err != nil {
			return nil, err
		}

		keywords = append(keywords, *keyword)
	}

	return &keywords, nil
}

func (s *PostgresStore) GetKeyword(id int) (*common.Keyword, error) {
	query := `SELECT * FROM keywords WHERE id = $1;`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	rows.Next()

	keyword, err := ScanIntoKeyword(rows)
	if err != nil {
		return nil, err
	}

	return keyword, nil
}

func (s *PostgresStore) AddKeyword(keyword *common.Keyword) (int, error) {
	query := `INSERT INTO keywords (guild_id, keyword, reaction) VALUES ($1, $2, $3);`

	_, err := s.db.Exec(query, keyword.GuildID, keyword.Keyword, keyword.Reaction)

	if err != nil {
		return -1, err
	}

	query = `SELECT id FROM keywords ORDER BY id DESC LIMIT 1;`

	row := s.db.QueryRow(query)

	var id int
	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *PostgresStore) UpdateKeyword(keyword *common.Keyword) error {
	query := `UPDATE keywords SET keyword = $1, reaction = $2 WHERE id = $3;`

	_, err := s.db.Exec(query, keyword.Keyword, keyword.Reaction, keyword.ID)
	return err
}

func (s *PostgresStore) DeleteKeyword(id int) error {
	query := `DELETE FROM keywords WHERE id = $1;`

	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) FindKeyword(guildID string, key string) (*common.Keyword, error) {
	query := `SELECT * FROM keywords WHERE guild_id = $1 AND keyword = $2;`

	rows, err := s.db.Query(query, guildID, key)
	if err != nil {
		return &common.Keyword{}, err
	}

	rows.Next()
	keyword, err := ScanIntoKeyword(rows)
	if err != nil {
		return &common.Keyword{}, err
	}

	return keyword, nil
}
