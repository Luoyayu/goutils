package main

import (
	"bytes"
	"context"
	"fmt"
	biliAPI "github.com/luoyayu/goutils/bilibili"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// params: video,sound,danmaku
func playLive(ctx context.Context, cid interface{}, params map[string]int, costumedMpvArgs string) {
	ret, err := biliAPI.GetRoomPlayUrl(cid)
	if err != nil {
		Logger.Error(err)
		return
	}

	if ret.Code != 0 || ret.Data == nil {
		Logger.Error(ret.Message)
		return
	}

	var streamUrl string

	for _, url := range ret.Data.Durl {
		resp, err := http.Head(url.Url)
		if err != nil || resp.StatusCode == 404 {
			continue
		}
		streamUrl = url.Url
		break
	}

	fmt.Println("play:", streamUrl)
	log.Printf("3 params: %+v\n", params)
	log.Printf("costumed mpv params: %+v\n", costumedMpvArgs)

	if streamUrl == "" {
		Logger.Error("no available live stream!")
	} else {
		if params["danmaku"] == 1 {
			go (&biliAPI.DanmakuClient{}).NewDanmakuClient(ctx, cid)
		}

		go func(ctx context.Context, cid interface{}, url string, params map[string]int, costumedMpvArgs string) {
			var playArgs []string

			// --cache=<yes|no|auto>
			// --no-cache
			// --cache --demuxer-max-bytes=10MB
			// --ontop

			if params["video"] == -1 && params["sound"] == -1 && len(costumedMpvArgs) == 0 {
				return
			}

			if params["video"] == -1 && params["sound"] == 1 && len(costumedMpvArgs) == 0 {
				playArgs = append(playArgs, "--no-video")
			}

			playArgs = append(playArgs, costumedMpvArgs, streamUrl)

			log.Println("mpv:", playArgs)

			cmd := exec.Command(mpv, playArgs...)
			out := &bytes.Buffer{}

			cmd.Stdout = out

			// stop mpv if exists
			stopMpvSafely()

			err = cmd.Start()

			if err == nil {
				MpvPid = int32(cmd.Process.Pid)
			}
			//log.Println("start err:", err)
			//log.Println("mpv pid:", cmd.Process.Pid)

			go func(out *bytes.Buffer) {
				tm := time.NewTicker(time.Second * 3)
				for {
					select {
					case <-tm.C:
						return
					default:
						if strings.Contains(out.String(), "option not found") {
							Logger.Error(strings.Split(out.String(), "\n")[0])
							return
						}
					}
				}
			}(out)

			select {
			case <-ctx.Done():
				//log.Println("context is canceled: kill mpv player")
				stopMpvSafely()
				return
			}

		}(ctx, cid, streamUrl, params, costumedMpvArgs)

		/*fmt.Println("video:", params["video"])
		fmt.Println("sound:", params["sound"])
		fmt.Println("danmaku:", params["danmaku"])
		fmt.Println("live url:", streamUrl)*/

		select {
		case <-ctx.Done():
			log.Println("received Stop signal")
			return
		}
	}
	log.Println("play Live END!")
}
