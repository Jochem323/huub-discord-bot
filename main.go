package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + "MTE2NDMyMTQ0MTQzMzUzODY5MA.G4zpAX.9w1SuTsQieI1vjO6tF2l2-JceVaJrqf37P_GQY")
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "pong")
		}
	})

	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}
