package main

import (
	"fmt"
	"github.com/luoyayu/netease_go/api"
	"log"
	"net/http"
	"strconv"
)

func download(w http.ResponseWriter, r *http.Request) {

	songId := r.URL.Query().Get("songid")

	log.Println("request singId: ", songId)
	songIdInt64, err := strconv.ParseInt(songId, 10, 64)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
	} else {
		s := api.GetSongDetail(songIdInt64)
		d := api.GetSongDownloadUrl(songId, -1)
		if d != nil && d.Code == 200 {
			_, _ = w.Write([]byte(
				fmt.Sprintf("%s </br>%s.%s", d.Data.Url, s.Name, d.Data.Type)),
			)
		} else {
			log.Println("error!")
		}
	}

}

func main() {
	http.HandleFunc("/download", download)
	log.Fatal(http.ListenAndServe("0.0.0.0:16333", nil))
}
