package storage

import (
	"database/sql"
	"huub-discord-bot/common"
)

func (s *PostgresStore) CreateguildTable() error {
	query := `CREATE TABLE IF NOT EXISTS guilds (
		id TEXT PRIMARY KEY,
		prefix TEXT NOT NULL
	);`

	_, err := s.db.Exec(query)
	return err
}

func ScanIntoGuild(rows *sql.Rows) (*common.Guild, error) {
	guild := new(common.Guild)
	err := rows.Scan(
		&guild.ID,
		&guild.Prefix,
	)

	return guild, err
}

func (s *PostgresStore) GetGuild(guildID string) (*common.Guild, error) {
	query := `SELECT * FROM guilds WHERE id = $1;`

	rows, err := s.db.Query(query, guildID)
	if err != nil {
		return nil, err
	}

	rows.Next()

	guild, err := ScanIntoGuild(rows)
	if err != nil {
		return nil, err
	}

	return guild, nil
}

func (s *PostgresStore) AddGuild(guild common.Guild) error {
	query := `INSERT INTO guilds (id, prefix) VALUES ($1, $2);`

	_, err := s.db.Exec(query, guild.ID, guild.Prefix)
	return err
}

func (s *PostgresStore) UpdateGuild(guild common.Guild) error {
	query := `UPDATE guilds SET prefix = $1 WHERE id = $2;`

	_, err := s.db.Exec(query, guild.Prefix, guild.ID)
	return err
}

func (s *PostgresStore) DeleteGuild(guildID string) error {
	query := `DELETE FROM guilds WHERE id = $1;`

	_, err := s.db.Exec(query, guildID)
	return err
}
