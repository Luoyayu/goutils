package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNetEaseClient_GetRecommend(t *testing.T) {
	assert.Equal(t, SINGED, true, "must login!")
	if songs, err := _client.GetRecommend(); err == nil {
		for i, song := range songs.Recommend {
			if len(song.Alias) == 0 {
				fmt.Printf("[%02d]: %s--%s\n", i+1, song.Name, song.Artists[0].Name)
			} else {
				//log.Println(song.Alias)
				alias := ""
				sep2 := strings.Split(song.Alias[0], "；")
				sep3 := strings.Split(song.Alias[0], "/")
				if len(sep2) == 2 {
					alias = strings.TrimSpace(sep2[0])
				} else if len(sep3) == 2 {
					alias = strings.TrimSpace(sep3[0])
				} else {
					alias = song.Alias[0]
				}
				fmt.Printf("[%02d]: %s--%s\n", i+1, song.Name+"("+alias+")", song.Artists[0].Name)
			}
		}
	} else {
		panic(err)
	}
}
