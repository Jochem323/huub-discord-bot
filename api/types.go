package api

import (
	"huub-discord-bot/common"
	"huub-discord-bot/storage"
	"log"
	"net/http"
)

type APIServer struct {
	ListenAdress string
	log          *log.Logger
	GuildStore   storage.GuildStore
	KeywordStore storage.KeywordStore
	APIKeyStore  storage.APIKeyStore
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type APIResponse struct {
	Message string `json:"message"`
}

type CreateGuildRequest struct {
	ID     string `json:"id"`
	Prefix string `json:"prefix"`
}

type UpdateGuildRequest struct {
	Prefix string `json:"prefix"`
}

type CreateKeywordRequest struct {
	GuildID  string `json:"guild_id"`
	Keyword  string `json:"keyword"`
	Reaction string `json:"reaction"`
}

type UpdateKeywordRequest struct {
	Reaction string `json:"reaction"`
}

type CreateKeyRequest struct {
	Admin     bool   `json:"admin"`
	GuildID   string `json:"guild_id"`
	Comment   string `json:"comment"`
	CreatedBy string `json:"created_by"`
}

type CreateKeyResponse struct {
	Key   common.APIKey `json:"key"`
	Token string        `json:"token"`
}

type UpdateKeyRequest struct {
	Comment   string `json:"comment"`
	Active    *bool  `json:"active"`
	Revoked   *bool  `json:"revoked"`
	Ratelimit *bool  `json:"ratelimit"`
}
