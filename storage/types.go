package storage

import (
	"database/sql"
	"huub-discord-bot/common"
)

type PostgresStore struct {
	db *sql.DB
}

type KeywordStore interface {
	// GetKeywords returns the keywords for a given guild
	GetKeywords(guildID string) ([]common.Keyword, error)

	// AddKeyword adds a keyword to the database
	AddKeyword(guildID string, keyword common.Keyword) error

	// UpdateKeyword updates a keyword in the database
	UpdateKeyword(keyword common.Keyword) error

	// DeleteKeyword deletes a keyword from the database
	DeleteKeyword(id int) error

	// FindKeyword finds a keyword in the database
	FindKeyword(guildID string, key string) (common.Keyword, error)
}

type GuildStore interface {
	// GetGuild returns the guild for a given guildID
	GetGuild(guildID string) (*common.Guild, error)

	// AddGuild adds a guild to the database
	AddGuild(guild common.Guild) error

	// UpdateGuild updates a guild in the database
	UpdateGuild(guild common.Guild) error

	// DeleteGuild deletes a guild from the database
	DeleteGuild(guildID string) error
}
