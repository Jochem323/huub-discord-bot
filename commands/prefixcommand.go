package commands

import (
	"fmt"
	"huub-discord-bot/storage"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func PrefixCommand(d *discordgo.Session, m *discordgo.MessageCreate, store storage.GuildStore) {
	// Check if the user is an admin
	if !IsAdmin(d, m) {
		d.ChannelMessageSend(m.ChannelID, "You are not an admin")
		return
	}

	substrings := strings.Split(m.Content, " ")
	if len(substrings) < 2 {
		d.ChannelMessageSend(m.ChannelID, "Please provide a prefix")
		return
	}

	if len(substrings) > 2 {
		d.ChannelMessageSend(m.ChannelID, "Prefix cannot contain spaces")
		return
	}

	// Retrieve the prefix from the message
	prefix := strings.TrimSpace(substrings[1])
	if prefix == "/" {
		d.ChannelMessageSend(m.ChannelID, "Prefix cannot be /")
		return
	}

	// Retrieve the guild from the database
	guild, err := store.GetGuild(m.GuildID)
	if err != nil {
		return
	}

	// Update the prefix in the database
	guild.Prefix = prefix
	err = store.UpdateGuild(guild)
	if err != nil {
		return
	}

	// Send a message to the channel
	d.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Prefix updated to %s", prefix))
}
