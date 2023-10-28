package storage

import (
	"database/sql"
	"huub-discord-bot/common"
)

func (s *PostgresStore) CreateGuildTable() error {
	query := `CREATE TABLE IF NOT EXISTS guilds (
		id TEXT PRIMARY KEY,
		prefix TEXT NOT NULL,
		api_enabled BOOLEAN NOT NULL DEFAULT FALSE
	);`

	_, err := s.db.Exec(query)
	return err
}

func ScanIntoGuild(rows *sql.Rows) (any, error) {
	guild := new(common.Guild)
	err := rows.Scan(
		&guild.ID,
		&guild.Prefix,
		&guild.APIEnabled,
	)

	return *guild, err
}

func (s *PostgresStore) GetGuilds() ([]common.Guild, error) {
	query := `SELECT * FROM guilds;`

	results, err := s.GetMultiple(ScanIntoGuild, query)
	if err != nil {
		return nil, err
	}

	guilds := []common.Guild{}
	for _, result := range results {
		guild := result.(common.Guild)
		guilds = append(guilds, guild)
	}

	return guilds, nil
}

func (s *PostgresStore) GetGuild(guildID string) (common.Guild, error) {
	query := `SELECT * FROM guilds WHERE id = $1;`

	result, err := s.GetOne(ScanIntoGuild, query, guildID)
	if err != nil {
		return common.Guild{}, err
	}

	guild := result.(common.Guild)

	return guild, nil
}

func (s *PostgresStore) AddGuild(guild common.Guild) error {
	query := `INSERT INTO guilds (id, prefix) VALUES ($1, $2);`
	return s.Exec(query, guild.ID, guild.Prefix)
}

func (s *PostgresStore) UpdateGuild(guild common.Guild) error {
	query := `UPDATE guilds SET prefix = $1, api_enabled = $2 WHERE id = $3;`
	return s.Exec(query, guild.Prefix, guild.APIEnabled, guild.ID)
}

func (s *PostgresStore) DeleteGuild(guildID string) error {
	query := `DELETE FROM guilds WHERE id = $1;`
	return s.Exec(query, guildID)
}
