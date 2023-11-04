package music

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (v *VCHandler) JoinCommand(d *discordgo.Session, guildID string, channelID string) error {
	_, ok := v.connections[guildID]
	if ok {
		return fmt.Errorf("already connected to a voice channel")
	}

	voice, err := d.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		return err
	}

	var queue []Song
	v.connections[guildID] = &vcConnection{
		guildID: guildID,
		voice:   voice,
		queue:   queue,
		stop:    make(chan bool),
	}

	return nil
}

func (v *VCHandler) SummonCommand(guildID string, channelID string) error {
	conn, ok := v.connections[guildID]
	if !ok {
		return fmt.Errorf("not connected to a voice channel")
	}

	err := conn.voice.ChangeChannel(channelID, false, false)
	if err != nil {
		return err
	}

	return nil
}

func (v *VCHandler) LeaveCommand(guildID string) error {
	conn, ok := v.connections[guildID]
	if !ok {
		return fmt.Errorf("not connected to a voice channel")
	}

	err := conn.voice.Disconnect()
	if err != nil {
		return err
	}

	delete(v.connections, guildID)

	return nil
}

func (v *VCHandler) PlayCommand(d *discordgo.Session, guildID string, channelID string, song Song) error {
	conn, ok := v.connections[guildID]
	if !ok {
		v.JoinCommand(d, guildID, channelID)
		conn, ok = v.connections[guildID]
		if !ok {
			return fmt.Errorf("error joining voice channel")
		}
	}

	v.log.Println("Adding song to queue:", song.URL)

	v.AddSongToQueue(conn, song)

	v.log.Println("Queue:", conn.queue)

	if conn.nowPlaying.URL == "" {
		v.log.Println("No song playing, playing next song")
		v.startQueue(conn)
	}

	return nil
}

func (v *VCHandler) StopCommand(d *discordgo.Session, guildID string) error {
	conn, ok := v.connections[guildID]
	if !ok {
		return fmt.Errorf("not connected to a voice channel")
	}

	conn.stop <- true

	return nil
}
