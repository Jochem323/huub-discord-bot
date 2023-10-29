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
	logger := log.New(os.Stdout, "Main: ", log.Ldate|log.Ltime)

	// Connect to the database
	db, err := storage.NewPostgresStore()
	if err != nil {
		logger.Fatal(err)
	}

	// Initialize the database
	err = db.Init()
	if err != nil {
		logger.Fatal(err)
	}

	// Initialize the Discord bot
	discord := discordbot.DiscordBot{
		KeywordStore: db,
		GuildStore:   db,
		APIKeyStore:  db,
	}
	err = discord.Init()
	if err != nil {
		logger.Fatal(err)
	}

	// Add handlers
	discord.AddHandler(discord.KeywordHandler)
	discord.AddHandler(discord.CommandHandler)

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

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}
