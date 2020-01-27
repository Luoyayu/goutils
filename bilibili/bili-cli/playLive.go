package main

import (
	"bytes"
	"context"
	biliAPI "github.com/luoyayu/goutils/bilibili"
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
		showPageLive()
		return
	}

	if ret.Code != 0 || ret.Data == nil {
		Logger.Error(ret.Message)
		showPageLive()
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

	//Logger.Info("play:", streamUrl)
	//Logger.Info("3 params: %v\n", params)
	//Logger.Infof("costumed mpv params: %s\n", costumedMpvArgs)

	if streamUrl == "" {
		Logger.Error("no available live stream!")
		showPageLive()
		return
	} else {
		if params["danmaku"] == 1 {
			go (&biliAPI.DanmakuClient{}).NewDanmakuClient(ctx, cid)
		}

		// FIXME: mpv more params can't parse!
		go func(ctx context.Context, cid interface{}, url string, params map[string]int, costumedMpvArgs string) {
			var playArgs = []string{}

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

			playArgs = append(playArgs, "--hwdec=auto")

			if costumedMpvArgs != "" {
				playArgs = append(playArgs, strings.Split(costumedMpvArgs, " ")...)
			}

			playArgs = append(playArgs, streamUrl)
			//Logger.Info("mpv:", playArgs)

			cmd := exec.Command(mpv, playArgs...)
			out := &bytes.Buffer{}
			cmd.Stdout = out

			stopMpvSafely()

			err = cmd.Start()

			if err == nil {
				MpvPid = int32(cmd.Process.Pid)
			}

			//Logger.Info("mpv pid:", cmd.Process.Pid)

			go func(out *bytes.Buffer) {
				tm := time.NewTicker(time.Second * 5)
				for {
					select {
					case <-tm.C:
						return
					default:
						//Logger.Info(out.String())
						if strings.Contains(out.String(), "option not found") {
							Logger.Error(strings.Split(out.String(), "\n"))
							return
						}
						time.Sleep(time.Second * 1)
					}
				}
			}(out)

			select {
			case <-ctx.Done():
				//Logger.Info("context is canceled: kill mpv player")
				stopMpvSafely()
				return
			}

		}(ctx, cid, streamUrl, params, costumedMpvArgs)

		select {
		case <-ctx.Done():
			//Logger.Info("received Stop signal")
			return
		}
	}
}
