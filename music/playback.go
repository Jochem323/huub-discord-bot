package music

import (
	"fmt"
	"time"

	"github.com/bwmarrin/dgvoice"
)

func (v *VCHandler) startQueue(conn *vcConnection) {
	go func() {
		for len(conn.queue) > 0 {
			err := v.playNextSong(conn)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
}

func (v *VCHandler) playNextSong(conn *vcConnection) error {
	if len(conn.queue) == 0 {
		return fmt.Errorf("no songs in queue")
	}

	song := conn.queue[0]
	conn.nowPlaying = song
	conn.queue = conn.queue[1:]

	v.log.Printf("Next up: %s", song.URL)

	for !song.Ready {
		time.Sleep(1 * time.Second)
	}

	v.log.Println("Playing song...")

	dgvoice.PlayAudioFile(conn.voice, song.Filepath, conn.stop)

	conn.nowPlaying = Song{}
	conn.stop <- false

	v.log.Println("End of playback")

	return nil
}
