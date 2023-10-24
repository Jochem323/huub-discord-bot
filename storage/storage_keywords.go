package storage

import (
	"database/sql"
	"huub-discord-bot/common"
)

func (s *PostgresStore) CreatekeywordTable() error {
	query := `CREATE TABLE IF NOT EXISTS keywords (
		id SERIAL PRIMARY KEY,
		guild_id TEXT NOT NULL,
		keyword TEXT NOT NULL,
		reaction TEXT NOT NULL
	);`

	_, err := s.db.Exec(query)
	return err
}

func ScanRowsIntoKeyword(rows *sql.Rows) (*common.Keyword, error) {
	keyword := new(common.Keyword)
	err := rows.Scan(
		&keyword.ID,
		&keyword.GuildID,
		&keyword.Key,
		&keyword.Reaction,
	)

	return keyword, err
}

func ScanRowIntoKeyword(row *sql.Row) (*common.Keyword, error) {
	keyword := new(common.Keyword)
	err := row.Scan(
		&keyword.ID,
		&keyword.GuildID,
		&keyword.Key,
		&keyword.Reaction,
	)

	return keyword, err
}

func (s *PostgresStore) GetKeywords(guildID string) ([]common.Keyword, error) {
	query := `SELECT * FROM keywords WHERE guild_id = $1;`

	rows, err := s.db.Query(query, guildID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keywords := []common.Keyword{}
	for rows.Next() {
		keyword, err := ScanRowsIntoKeyword(rows)
		if err != nil {
			return nil, err
		}

		keywords = append(keywords, *keyword)
	}

	return keywords, nil
}

func (s *PostgresStore) AddKeyword(guildID string, keyword common.Keyword) error {
	query := `INSERT INTO keywords (guild_id, keyword, reaction) VALUES ($1, $2, $3);`

	_, err := s.db.Exec(query, guildID, keyword.Key, keyword.Reaction)
	return err
}

func (s *PostgresStore) UpdateKeyword(keyword common.Keyword) error {
	query := `UPDATE keywords SET keyword = $1, reaction = $2 WHERE id = $3;`

	_, err := s.db.Exec(query, keyword.Key, keyword.Reaction, keyword.ID)
	return err
}

func (s *PostgresStore) DeleteKeyword(id int) error {
	query := `DELETE FROM keywords WHERE id = $1;`

	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) FindKeyword(guildID string, key string) (common.Keyword, error) {
	query := `SELECT * FROM keywords WHERE guild_id = $1 AND keyword = $2;`

	row := s.db.QueryRow(query, guildID, key)

	keyword, err := ScanRowIntoKeyword(row)
	if err != nil {
		return common.Keyword{}, err
	}

	return *keyword, nil
}
