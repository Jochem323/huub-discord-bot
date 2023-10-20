package commands

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

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

func MonkeCommand(d *discordgo.Session, m *discordgo.MessageCreate) {
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

	d.ChannelMessageSend(m.ChannelID, gifResponse.Results[0].Media.Gif.URL)
}
