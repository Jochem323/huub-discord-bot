package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"huub-discord-bot/api"
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

	// Get the listen address from the environment
	listenAddress, found := os.LookupEnv("API_LISTEN_ADDRESS")
	if !found {
		// Default for running in Docker
		listenAddress = "0.0.0.0:8080"
	}

	// Initialize the API server
	api := api.APIServer{
		ListenAdress: listenAddress,
		KeywordStore: db,
		GuildStore:   db,
		APIKeyStore:  db,
	}
	api.Run()

	// Initialize the Discord bot
	discord := discordbot.NewDiscordBot(db, db)
	err = discord.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Add handlers
	discord.AddHandler(discord.KeywordHandler)
	discord.AddHandler(discord.CommandHandler)

	// Wait for a signal to close the bot
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}
