package discordbot

import (
	"huub-discord-bot/storage"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	session      *discordgo.Session
	keywordStore storage.KeywordStore
	guildStore   storage.GuildStore
}

func NewDiscordBot(keywordStore storage.KeywordStore, guildStore storage.GuildStore) *DiscordBot {
	return &DiscordBot{
		keywordStore: keywordStore,
		guildStore:   guildStore,
	}
}
