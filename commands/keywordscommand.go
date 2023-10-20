package commands

import (
	"huub-discord-bot/common"
	"huub-discord-bot/storage"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func KeywordsCommand(d *discordgo.Session, m *discordgo.MessageCreate, keywordStore storage.KeywordStore) {
	// Get the subcommand
	subcommand := strings.Split(m.Content, " ")[1]

	switch subcommand {
	case "add":
		AddKeyword(d, m, keywordStore)
	case "update":
		UpdateKeyword(d, m, keywordStore)
	case "remove":
		RemoveKeyword(d, m, keywordStore)
	}
}

func AddKeyword(d *discordgo.Session, m *discordgo.MessageCreate, keywordStore storage.KeywordStore) {
	// Get the keyword
	key := strings.Split(m.Content, " ")[2]

	// Get the reaction
	reaction := strings.Split(m.Content, " ")[3]

	// Get the guildID
	guildID := m.GuildID

	// Create the keyword
	keyword := common.NewKeyword(guildID, key, reaction)

	// Create the keyword
	err := keywordStore.AddKeyword(guildID, keyword)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	d.ChannelMessageSend(m.ChannelID, "Keyword added")
}

func UpdateKeyword(d *discordgo.Session, m *discordgo.MessageCreate, keywordStore storage.KeywordStore) {
	// Get the keyword
	key := strings.Split(m.Content, " ")[2]

	// Get the reaction
	reaction := strings.Split(m.Content, " ")[3]

	// Get the guildID
	guildID := m.GuildID

	// Get the old keyword
	keyword, err := keywordStore.FindKeyword(guildID, key)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	// Create the new keyword
	newKeyword := common.NewKeyword(guildID, key, reaction)
	newKeyword.ID = keyword.ID

	// Create the keyword
	err = keywordStore.UpdateKeyword(newKeyword)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	d.ChannelMessageSend(m.ChannelID, "Keyword updated")
}

func RemoveKeyword(d *discordgo.Session, m *discordgo.MessageCreate, keywordStore storage.KeywordStore) {
	// Get the keyword
	key := strings.Split(m.Content, " ")[2]

	// Get the guildID
	guildID := m.GuildID

	keyword, err := keywordStore.FindKeyword(guildID, key)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	// Delete the keyword
	err = keywordStore.DeleteKeyword(keyword.ID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	d.ChannelMessageSend(m.ChannelID, "Keyword deleted")
}
