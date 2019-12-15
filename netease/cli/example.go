package main

import (
	"fmt"
	"github.com/luoyayu/goutils/netease/api"
	"github.com/spf13/afero"
	"log"
	"os"
	"strings"
)

var client = &api.NetEaseClient{}

func init() {
	client.NewDefaultHttpClient()
}

func exampleDownloadSong() {
	song := api.SongStruct{Id: 560079}
	if err := song.Download(client, "", -1); err != nil {
		panic(err)
	}
}

func exampleDownloadSongs() {
	songs := api.SongsStruct{
		Songs: []api.SongStruct{
			{Id: 560079},
			{Id: 1387099649},
			{Id: 514774419},
		},
	}

	if err := songs.Download(client, "", nil); err != nil {
		panic(err)
	}
}

func exampleGetSongDetail() {
	song := client.GetSongDetail(560079)
	fmt.Printf("%+v\n", song)
}

func exampleLoginCellphone() {
	_ = afero.NewOsFs().Remove("cookies")
	if err := client.LoginByCellphone(os.Getenv("PHONE"), os.Getenv("PASSWD")); err != nil {
		log.Println(err)
	}
}

func exampleSaveCookies() {

	if err := api.SaveCookies("cookies_gob", map[string]string{"name1": "1", "name2": "22222222",}); err != nil {
		panic(err)
	} else {
		log.Println("save !")
	}
}

func exampleLoadCookies() {
	if ret, err := api.LoadCookies("cookies_gob"); err != nil {
		panic(err)
	} else {
		log.Println("read from cookies_gob", ret)

	}
}

func exampleGetRecommendSong() {
	if songs, err := client.GetRecommend(); err == nil {
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

func main() {

	//exampleDownloadSong()
	//exampleDownloadSongs()
	//exampleGetSongDetail()

	//exampleLoginCellphone()
	//exampleSaveCookies()
	//exampleLoadCookies()

	exampleGetRecommendSong()
}
