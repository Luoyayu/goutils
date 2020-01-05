package api

import (
	"encoding/json"
	"errors"
	"github.com/spf13/afero"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var appFs afero.Fs

func init() {
	appFs = afero.NewOsFs()
}

type DownloadDataStruct struct {
	Data struct {
		Id    int64  `json:"id"`
		Url   string `json:"url"`
		Size  int64  `json:"size"`
		Md5   string `json:"md5"`
		Type  string `json:"type"`
		Fee   int64  `json:"fee"`
		Payed int64  `json:"payed"`
		Flag  int    `json:"flag"`
	} `json:"data"`
	Code int64 `json:"code"`
}

const (
	DownloadRAW  = 999000
	DownloadHigh = 320000
	DownloadMid  = 160000
	DownloadLow  = 96000
)

// GetSongDownloadUrl 指定比特率, 从歌曲信息中获取, 实际上限更高, 非正数表示最高比特率
func (r *NetEaseClient) GetSongDownloadUrl(songId string, br int64) *DownloadDataStruct {
	// br: 999000, 320000, 160000, 96000
	if br <= 0 {
		br = DownloadRAW
	}

	params := map[string]interface{}{
		"e_r": false,
		"id":  songId,
		"br":  br,
	}
	_ = r.NewEnc("/api/song/enhance/download/url", params)
	if ret, err := r.DoPost(); err == nil {
		downloadInfo := DownloadDataStruct{}
		err = json.Unmarshal(ret, &downloadInfo)
		if downloadInfo.Code == 200 {
			debugV("downloadInfo", downloadInfo)
			return &downloadInfo
		}
	} else {
		log.Println(err)
	}
	return nil
}

// DownloadSong : download a song
func DownloadSong(songName, songType, songUrl string, path string) (err error) {
	if path == "" {
		path, _ = os.Getwd()
	}
	if songUrl == "" {
		return errors.New("song url is None")
	}

	if HeaderMap["MUSIC_U"] == "" {
		return errors.New("please login first!")
	}

	path = strings.Join([]string{path, songName + "." + songType}, string(os.PathSeparator))
	var exist bool
	var f afero.File
	var req *http.Request
	var resp *http.Response

	if exist, err = afero.Exists(appFs, path); err == nil && exist == false {
		if f, err = appFs.Create(path); err == nil {
			c := &http.Client{}
			if req, err = http.NewRequest("GET", songUrl, nil); err == nil {
				req.Header.Set("host", "music.163.com")
				req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/605.1.15 (KHTML, like Gecko)")
				req.Header.Set("cookies", "channel=netease")
				if resp, err = c.Do(req); err == nil && resp != nil && resp.Body != nil {
					_, _ = io.Copy(f, resp.Body)
					log.Println("finish download: ", path)
					defer f.Close()
					defer resp.Body.Close()
				}
			}
		}
	} else {
		err = errors.New(path + ": file already exists")
	}
	return
}

// DownloadSongs:
func DownloadSongs(songsName, songsType, songsUrl []string, path string) (err error) {
	size := len(songsName)
	if size != len(songsType) || size != len(songsUrl) {
		return errors.New("songs Name, songs Type, songs Url not match")
	}

	var wg sync.WaitGroup
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(idx int) {
			err := DownloadSong(songsName[i], songsType[i], songsUrl[i], path)
			if err != nil {
				log.Println(songsName[i], "download failed: ", err)
			}
			wg.Done()
		}(i)
		time.Sleep(time.Microsecond * 100)
	}
	wg.Wait()
	return err
}
