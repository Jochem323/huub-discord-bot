package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type GifResponse struct {
	Results []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Media struct {
			Gif struct {
				URL string `json:"url"`
			} `json:"gif"`
		} `json:"media_formats"`
	} `json:"results"`
}

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
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

		if m.Content == "PING" {
			s.ChannelMessageSend(m.ChannelID, "PONG")
		}

		if m.Content == "monke" {
			url := "https://tenor.googleapis.com/v2/search?q=monke&key=" + os.Getenv("TENOR_KEY") + "&client_key=huub-discord-bot&country=NL&random=true&limit=50&media_filter=gif"
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("client: could not read response body: %s\n", err)
			}

			var gifResponse GifResponse
			err = json.Unmarshal(respBody, &gifResponse)
			if err != nil {
				log.Println(err)
			}

			s.ChannelMessageSend(m.ChannelID, gifResponse.Results[0].Media.Gif.URL)
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
