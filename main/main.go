package main

import (
	"log"
	"os"

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

	// Initialize the Discord bot
	discord := discordbot.DiscordBot{
		KeywordStore: db,
		GuildStore:   db,
		APIKeyStore:  db,
	}
	err = discord.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Add handlers
	discord.AddHandler(discord.KeywordHandler)
	discord.AddHandler(discord.CommandHandler)

	log.Println("Discord bot running")

	defer discord.Close()

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
}
