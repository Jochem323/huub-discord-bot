package discordbot

import (
	"huub-discord-bot/storage"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	session      *discordgo.Session
	KeywordStore storage.KeywordStore
	GuildStore   storage.GuildStore
	APIKeyStore  storage.APIKeyStore
}
