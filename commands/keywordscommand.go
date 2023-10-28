package commands

import (
	"huub-discord-bot/common"
	"huub-discord-bot/storage"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func KeywordsCommand(d *discordgo.Session, m *discordgo.MessageCreate, keywordStore storage.KeywordStore) {
	// Check if the user is an admin
	if !IsAdmin(d, m) {
		d.ChannelMessageSend(m.ChannelID, "You are not an admin")
		return
	}

	// Get the subcommand
	subcommand := strings.Split(m.Content, " ")[1]

	switch subcommand {
	case "add":
		AddKeyword(d, m, keywordStore)
	case "list":
		ListKeywords(d, m, keywordStore)
	case "update":
		UpdateKeyword(d, m, keywordStore)
	case "remove":
		RemoveKeyword(d, m, keywordStore)
	}
}

func AddKeyword(d *discordgo.Session, m *discordgo.MessageCreate, keywordStore storage.KeywordStore) {
	// Get substrings
	substrings := strings.SplitN(m.Content, " ", 4)

	// Check if the substrings are correct
	if len(substrings) != 4 {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	// Create the keyword
	keyword := common.Keyword{
		GuildID:  m.GuildID,
		Keyword:  substrings[2],
		Reaction: substrings[3],
	}

	// Create the keyword
	_, err := keywordStore.AddKeyword(keyword)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	d.ChannelMessageSend(m.ChannelID, "Keyword added")
}

func ListKeywords(d *discordgo.Session, m *discordgo.MessageCreate, keywordStore storage.KeywordStore) {
	// Get the guildID
	guildID := m.GuildID

	// Get the keywords
	keywords, err := keywordStore.GetKeywords(guildID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	// Create the message
	message := "Keywords:\n"
	for _, keyword := range keywords {
		message += keyword.Keyword + ": " + keyword.Reaction + "\n"
	}

	d.ChannelMessageSend(m.ChannelID, message)
}

func UpdateKeyword(d *discordgo.Session, m *discordgo.MessageCreate, keywordStore storage.KeywordStore) {
	// Get substrings
	substrings := strings.SplitN(m.Content, " ", 4)

	// Check if the substrings are correct
	if len(substrings) != 4 {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	// Get the old keyword
	keyword, err := keywordStore.FindKeyword(m.GuildID, substrings[2])
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	keyword.Reaction = substrings[3]

	// Create the keyword
	err = keywordStore.UpdateKeyword(keyword)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Something went wrong")
		return
	}

	d.ChannelMessageSend(m.ChannelID, "Keyword updated")
}

func RemoveKeyword(d *discordgo.Session, m *discordgo.MessageCreate, keywordStore storage.KeywordStore) {
	// Get the keyword
	key := strings.Split(m.Content, " ")[2]

	keyword, err := keywordStore.FindKeyword(m.GuildID, key)
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
