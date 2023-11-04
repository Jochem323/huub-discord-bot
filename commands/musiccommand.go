package commands

import (
	"huub-discord-bot/music"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func MusicCommand(d *discordgo.Session, m *discordgo.MessageCreate, v *music.VCHandler) {
	subcommand := strings.Split(m.Content, " ")[1]

	switch subcommand {
	case "join":
		JoinCommand(d, m, v)
	case "summon":
		SummonCommand(d, m, v)
	case "leave":
		LeaveCommand(d, m, v)
	case "play":
		PlayCommand(d, m, v)
	case "stop":
		StopCommand(d, m, v)
	}
}

func JoinCommand(d *discordgo.Session, m *discordgo.MessageCreate, v *music.VCHandler) {
	guild, err := d.State.Guild(m.GuildID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Error getting guild")
		return
	}

	channelID := ""

	voiceStates := guild.VoiceStates
	for _, state := range voiceStates {
		if state.UserID == m.Author.ID {
			channelID = state.ChannelID
		}
	}

	if channelID == "" {
		d.ChannelMessageSend(m.ChannelID, "You are not in a voice channel")
	}

	err = v.JoinCommand(d, guild.ID, channelID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Error joining voice channel")
	}

	// dgvoice.PlayAudioFile(voice, "/home/jochem/nieuw.opus", make(chan bool))
}

func SummonCommand(d *discordgo.Session, m *discordgo.MessageCreate, v *music.VCHandler) {
	guild, err := d.State.Guild(m.GuildID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Error getting guild")
		return
	}

	channelID := ""

	voiceStates := guild.VoiceStates
	for _, state := range voiceStates {
		if state.UserID == m.Author.ID {
			channelID = state.ChannelID
		}
	}

	if channelID == "" {
		d.ChannelMessageSend(m.ChannelID, "You are not in a voice channel")
	}

	err = v.SummonCommand(guild.ID, channelID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Error changing voice channel")
	}
}

func LeaveCommand(d *discordgo.Session, m *discordgo.MessageCreate, v *music.VCHandler) {
	err := v.LeaveCommand(m.GuildID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Error leaving voice channel")
	}
}

func PlayCommand(d *discordgo.Session, m *discordgo.MessageCreate, v *music.VCHandler) {
	guild, err := d.State.Guild(m.GuildID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Error getting guild")
		return
	}

	channelID := ""

	voiceStates := guild.VoiceStates
	for _, state := range voiceStates {
		if state.UserID == m.Author.ID {
			channelID = state.ChannelID
		}
	}

	if channelID == "" {
		d.ChannelMessageSend(m.ChannelID, "You are not in a voice channel")
	}

	err = v.PlayCommand(d, guild.ID, channelID, music.Song{
		URL:         "https://www.youtube.com/watch?v=5qap5aO4i9A",
		Filepath:    "/home/jochem/nieuw.opus",
		User:        m.Author.ID,
		RequestTime: time.Now(),
		Ready:       true,
	})
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Error playing song")
	}
}

func StopCommand(d *discordgo.Session, m *discordgo.MessageCreate, v *music.VCHandler) {
	err := v.StopCommand(d, m.GuildID)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, "Error stopping song")
	}
}
