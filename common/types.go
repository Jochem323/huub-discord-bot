package common

import "time"

type Guild struct {
	ID     string `json:"id"`
	Prefix string `json:"prefix"`
}

type Keyword struct {
	ID       int    `json:"id"`
	GuildID  string `json:"guild_id"`
	Keyword  string `json:"keyword"`
	Reaction string `json:"reaction"`
}

type APIKey struct {
	ID        int       `json:"id"`
	Admin     bool      `json:"admin"`
	GuildID   string    `json:"guild_id"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
	Revoked   bool      `json:"revoked"`
	Ratelimit bool      `json:"ratelimit"`
}
