package music

import (
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type VCHandler struct {
	log         *log.Logger
	connections map[string]*vcConnection
}

func NewVCHandler() *VCHandler {
	handler := &VCHandler{
		log:         log.New(os.Stdout, "VCHandler ", log.Ldate|log.Ltime),
		connections: make(map[string]*vcConnection),
	}

	return handler
}

type vcConnection struct {
	guildID    string
	voice      *discordgo.VoiceConnection
	queue      []Song
	nowPlaying Song
	stop       chan bool
}

type Song struct {
	URL         string
	Filepath    string
	User        string
	RequestTime time.Time
	Ready       bool
}
