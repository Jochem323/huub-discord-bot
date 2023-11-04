package discordbot

import (
	"huub-discord-bot/music"
	"huub-discord-bot/storage"
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	log          *log.Logger
	session      *discordgo.Session
	VCHandler    *music.VCHandler
	KeywordStore storage.KeywordStore
	GuildStore   storage.GuildStore
	APIKeyStore  storage.APIKeyStore
}
