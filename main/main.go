package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"huub-discord-bot/discordbot"
	"huub-discord-bot/storage"
)

func main() {
	// Connect to the database
	db, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the database
	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}

	discord := discordbot.NewDiscordBot(db, db)
	err = discord.Init()
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(discord.KeywordHandler)
	discord.AddHandler(discord.CommandHandler)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}
