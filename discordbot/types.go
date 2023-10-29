package discordbot

import (
	"huub-discord-bot/storage"
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	log          *log.Logger
	session      *discordgo.Session
	KeywordStore storage.KeywordStore
	GuildStore   storage.GuildStore
	APIKeyStore  storage.APIKeyStore
}
