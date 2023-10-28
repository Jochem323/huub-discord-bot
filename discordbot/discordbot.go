package discordbot

import (
	"huub-discord-bot/commands"
	"huub-discord-bot/common"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *DiscordBot) Init() error {
	// Create a new Discord session using the provided bot token
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		return err
	}

	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = discord.Open()
	if err != nil {
		return err
	}

	b.session = discord

	b.AddGuildsToDB()

	log.Println("Bot is running")

	return nil
}

func (b *DiscordBot) Close() error {
	return b.session.Close()
}

func (b *DiscordBot) AddHandler(handler interface{}) {
	b.session.AddHandler(handler)
}

func (b *DiscordBot) KeywordHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Filter out messages from the bot itself and other bots
	if m.Author.ID == b.session.State.User.ID || m.Author.Bot {
		return
	}

	// Filter out messages that are not sent in a guild (so DMs)
	if m.GuildID == "" {
		return
	}

	keyword, err := b.KeywordStore.FindKeyword(m.GuildID, m.Content)
	if err != nil {
		return
	}

	s.ChannelMessageSend(m.ChannelID, keyword.Reaction)
}

func (b *DiscordBot) CommandHandler(d *discordgo.Session, m *discordgo.MessageCreate) {
	// Filter out messages from the bot itself and other bots
	if m.Author.ID == d.State.User.ID || m.Author.Bot {
		return
	}

	// Filter out messages that are not sent in a guild (so DMs)
	if m.GuildID == "" {
		return
	}

	// Get the guild from the database
	guild, err := b.GuildStore.GetGuild(m.GuildID)
	if err != nil {
		return
	}

	// Filter out messages that do not start with the command prefix
	if !strings.HasPrefix(m.Content, guild.Prefix) {
		return
	}

	// Retrieve the command word from the message
	command := strings.ToLower(strings.Split(m.Content, " ")[0][len(guild.Prefix):])

	switch command {
	case "ping":
		commands.PingCommand(d, m)
	case "monke":
		commands.MonkeCommand(d, m)
	case "keywords":
		commands.KeywordsCommand(d, m, b.KeywordStore)
	case "prefix":
		commands.PrefixCommand(d, m, b.GuildStore)
	case "api":
		commands.APICommand(d, m, b.GuildStore, b.APIKeyStore)
	}
}

func (d *DiscordBot) AddGuildsToDB() {
	guilds := d.session.State.Guilds

	for _, guild := range guilds {
		_, err := d.GuildStore.GetGuild(guild.ID)
		if err != nil {
			newGuild := common.Guild{
				ID:     guild.ID,
				Prefix: ".",
			}
			d.GuildStore.AddGuild(newGuild)
		}
	}
}
