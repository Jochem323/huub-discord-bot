package commands

import (
	"huub-discord-bot/api"
	"huub-discord-bot/common"
	"huub-discord-bot/storage"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func APICommand(d *discordgo.Session, m *discordgo.MessageCreate, guildStore storage.GuildStore, keyStore storage.APIKeyStore) {
	if !IsAdmin(d, m) {
		d.ChannelMessageSend(m.ChannelID, "You are not an admin")
		return
	}

	substrings := strings.Split(m.Content, " ")
	if len(substrings) < 2 {
		d.ChannelMessageSend(m.ChannelID, "Invalid command")
		return
	}

	subcommand := substrings[1]

	switch subcommand {
	case "key":
		APIKeyCommand(d, m, keyStore)
	case "adminkey":
		AdminAPIKeyCommand(d, m, keyStore)
	case "enable":
		EnableAPI(d, m, guildStore)
	case "disable":
		DisableAPI(d, m, guildStore)
	}
}

func EnableAPI(d *discordgo.Session, m *discordgo.MessageCreate, guildStore storage.GuildStore) {
	guild, err := guildStore.GetGuild(m.GuildID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "An error occured")
		return
	}

	if !guild.APIEnabled {
		guild.APIEnabled = true
		err := guildStore.UpdateGuild(guild)
		if err != nil {
			log.Println(err)
			d.ChannelMessageSend(m.ChannelID, "An error occured")
			return
		}
		d.ChannelMessageSend(m.ChannelID, "API enabled")
	} else {
		d.ChannelMessageSend(m.ChannelID, "API already enabled")
	}
}

func DisableAPI(d *discordgo.Session, m *discordgo.MessageCreate, guildStore storage.GuildStore) {
	guild, err := guildStore.GetGuild(m.GuildID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "An error occured")
		return
	}

	if guild.APIEnabled {
		guild.APIEnabled = false
		guildStore.UpdateGuild(guild)
		d.ChannelMessageSend(m.ChannelID, "API disabled")
	} else {
		d.ChannelMessageSend(m.ChannelID, "API already disabled")
	}
}

func APIKeyCommand(d *discordgo.Session, m *discordgo.MessageCreate, keyStore storage.APIKeyStore) {
	substrings := strings.Split(m.Content, " ")
	if len(substrings) < 3 {
		d.ChannelMessageSend(m.ChannelID, "Invalid command")
		return
	}
	subcommand := substrings[2]

	switch subcommand {
	case "add":
		AddAPIKey(false, d, m, keyStore)
	case "list":
		// ListAPIKeys(d, m)
	case "update":
		// UpdateAPIKey(d, m)
	case "remove":
		// RemoveAPIKey(d, m)
	}
}

func AdminAPIKeyCommand(d *discordgo.Session, m *discordgo.MessageCreate, keyStore storage.APIKeyStore) {
	substrings := strings.Split(m.Content, " ")
	if len(substrings) < 4 {
		d.ChannelMessageSend(m.ChannelID, "Invalid command")
		return
	}
	subcommand := substrings[2]

	pass := os.Getenv("BOT_TOKEN")[len(os.Getenv("BOT_TOKEN"))-10:]

	if substrings[3] != pass {
		d.ChannelMessageSend(m.ChannelID, "Invalid password")
		return
	}

	switch subcommand {
	case "add":
		AddAPIKey(true, d, m, keyStore)
	case "list":
		// ListAPIKeys(d, m)
	case "update":
		// UpdateAPIKey(d, m)
	case "remove":
		// RemoveAPIKey(d, m)
	}
}

func AddAPIKey(admin bool, d *discordgo.Session, m *discordgo.MessageCreate, keyStore storage.APIKeyStore) {
	substrings := strings.SplitN(m.Content, " ", 4)

	comment := ""
	if len(substrings) == 4 && !admin {
		comment = substrings[3]
	}

	guildId := ""
	if !admin {
		guildId = m.GuildID
	}

	key := common.APIKey{
		Admin:     admin,
		GuildID:   guildId,
		Comment:   comment,
		CreatedBy: m.Author.ID,
		CreatedAt: m.Timestamp,
		Active:    true,
		Revoked:   false,
		Ratelimit: true,
	}

	id, err := keyStore.AddKey(key)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	key.ID = id

	jwt, err := api.GenerateJWT(&key)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "An error occured")
		return
	}

	d.ChannelMessageSend(m.ChannelID, "Key added: "+jwt)
}
