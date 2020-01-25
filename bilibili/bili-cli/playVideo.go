package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func playVideo(ctx context.Context, aid interface{}, durl string, title string) {

	// --http-header-fields='Field1: value1','Field2: value2'
	header := strings.Join([]string{
		`Cookie: SESSDATA: ` + AccountSelected.SESSDATA,
		`User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:56.0) Gecko/20100101 Firefox/56.0`,
		`Referer: https://api.bilibili.com/x/web-interface/view?aid=` + fmt.Sprint(aid),
		`Origin: https://www.bilibili.com`,
		`bytes: 0-`,
	}, ",")
	header = `--http-header-fields=` + header
	var playArgs []string
	// HACK for osx

	playArgs = []string{header, "--force-media-title=" + title, "--hwdec=auto-safe" + durl}

	cmd := exec.Command(mpv, playArgs...)

	out := &bytes.Buffer{}
	cmd.Stdout = out

	err := cmd.Start()

	if err == nil {
		MpvPid = int32(cmd.Process.Pid)
	}

	//Logger.Info("mpv pid:", cmd.Process.Pid)

	// check
	/*go func(out *bytes.Buffer) {
		tm := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-tm.C:
				return
			default:
				Logger.Info(out.String())
				time.Sleep(time.Second * 1)
			}
		}
	}(out)*/

	select {
	case <-ctx.Done():
		stopMpvSafely()
		return
	}
}
