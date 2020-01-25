package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	biliAPI "github.com/luoyayu/goutils/bilibili"
	"strconv"
	"strings"
)

var NowLive []*Live

// for PlaySelected CMDBack
var pageLivePlaySelected2Select bool
var pageLivePlaySelected2Add bool
var pageLivePlaySelected2Recommend bool

func pageLivePlaySelectedBack() {
	if pageLivePlaySelected2Recommend { // Recommend
		pageLivePlaySelected2Recommend = false
		pageLiveRecommend()
	} else if pageLivePlaySelected2Select { // Select
		pageLivePlaySelected2Select = false
		pageLiveSelect(true)
	} else if pageLivePlaySelected2Add { // Add
		pageLivePlaySelected2Add = true
		pageLiveAdd()
	}
}

// end for PlaySelected CMDBack

func updateNowLive() {
	NowLive = pageSyncCmdSyncLive(false, false, false)
}

// if updateAccountNowLive is false
// you need to update NowLive []*Live manually
func pageLiveSelect(updateAccountNowLive bool) {
	if updateAccountNowLive {
		updateNowLive()
	}

	liveSelectableItems := make([]string, 0, len(NowLive))
	for _, ll := range NowLive {
		if ll.State == 1 && ll.Blocked == 0 {
			liveSelectableItems = append(liveSelectableItems, fmt.Sprintf("[%d] %s: %s", ll.Cid, ll.NikeName, ll.Title))
		}
	}

	// 跳转：nobody online -> showPageLive
	if len(liveSelectableItems) == 0 {
		Logger.Warn("nobody online!")
		showPageLive()
		return
	}

	result := strings.Split(PromptSelect("online live", liveSelectableItems, survey.WithPageSize(10)), " ")

	if len(result) > 0 && len(result[0]) > 3 {
		ridStr := result[0][1 : len(result[0])-1]

		if pageLivePlaySelected2Recommend == true {

		} else {
			pageLivePlaySelected2Select = true
		}
		pageLivePlaySelected(strings.TrimSpace(ridStr)) // 跳转：Select -> PlaySelected(rid)
	} else {
		showPageLive() // 跳转：Ctrl^C -> showPageLive
		return
	}
}

func pageLivePlaySelected(cid interface{}) {
	r, err := biliAPI.RoomInit(cid)
	if err != nil {
		Logger.Error("get room info failed!")
		showPageLive()
		return
	}

	u, err := biliAPI.GetUserInfo(r.Data.Uid)
	if err != nil {
		Logger.Error("get room host info failed!")
		showPageLive()
		return
	}

	uName := ""
	if u.Code == 0 {
		uName = u.Data.Name
	} else {
		uName = fmt.Sprint(r.Data.Uid)
	}

	if r.Code != 0 {
		Logger.Error(r.Message)
		showPageLive()
		return
	} else if r.Data.LiveStatus == 0 {
		Logger.Error("the live room for ", uName, " now is closed!")
		showPageLive()
		return
	}

	playLiveSelectedItems := []string{
		"play: video,sound,danmaku",
		"play: sound",
		"play: video,sound",
		"play: sound,danmaku",
		"play: danmaku",
		"play: costumed mpv args",
		CMDBack,
		CMDHome,
		CMDExit,
	}

	playCmd := PromptSelect("play live for "+uName, playLiveSelectedItems, survey.WithPageSize(len(playLiveSelectedItems)))
	//Logger.Info("playCmd: ", playCmd)

	if strings.HasPrefix(playCmd, "play: ") == false {
		switch playCmd {
		case CMDBack, "": // Ctrl^C or CMDBack
			pageLivePlaySelectedBack()
			return
		case CMDHome:
			showPageHome()
			return
		case CMDExit:
			exitClear()
		}
	}

	if strings.HasPrefix(playCmd, "play:") == true {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		paramsMap := map[string]int{
			"danmaku": -1,
			"video":   -1,
			"sound":   -1,
		} // 1:on, -1:off, 0:not set
		mpvOptions := ""

		if strings.Contains(playCmd, "mpv") {
			mpvOptions = PromptInput(&survey.Input{
				Message: "mpv and danmaku options",
				Default: "--danmaku=yes",
				Help:    "mpv manual: mpv.io/manual/stable \n recommend options:\n--ontop: keep player on top\n",
			})

			if strings.Contains(mpvOptions, "--danmaku=yes") {
				paramsMap["danmaku"] = 1
				mpvOptions = strings.Replace(mpvOptions, "--danmaku=yes", "", -1)
			} else if strings.Contains(mpvOptions, "--danmaku=no") {
				paramsMap["danmaku"] = -1
				mpvOptions = strings.Replace(mpvOptions, "--danmaku=no", "", -1)
			}

		} else {
			params := strings.Split(strings.TrimPrefix(playCmd, "play: "), ",")
			for i := range params {
				paramsMap[params[i]] = 1
			}
		}

		go playLive(ctx, cid, paramsMap, mpvOptions)

		pageLivePlayControl(cid, playCmd, cancel)
	}
}

func pageLivePlayControl(cid interface{}, playCmd string, cancel context.CancelFunc) {
	controlCallback := PromptSelect(
		"control live with "+playCmd,
		[]string{
			CMDBack,
			CMDHome,
			CMDExit,
		},
	)

	switch controlCallback {
	case CMDBack:
		cancel()
		pageLivePlaySelected(cid)
		return
	case CMDHome:
		cancel()
		showPageHome()
	case CMDExit:
		cancel()
		exitClear()
	}
}

func pageLiveAdd() {
	input := PromptInput(&survey.Input{
		Message: "rid:XXX / uid:XXX",
		Default: "",
		Help:    "add to play，input format must be `rid:`roomID or `uid:`up's uid",
	}, survey.WithValidator(
		func(ans interface{}) error {
			s := fmt.Sprint(ans)
			if strings.HasPrefix(s, "rid:") == false && strings.HasPrefix(s, "uid:") == false {
				return errors.New("input must start with `rid:` or `uid:`")
			} else {
				sp := strings.Split(s, ":")
				if len(sp) != 2 {
					return errors.New("wrong format, input must be like `rid:XXX` or `uid:YYY`")
				}
				if _, err := strconv.ParseInt(sp[1], 10, 64); err != nil {
					return errors.New("input must end with a number")
				}
			}
			return nil
		}),
	)

	sp := strings.Split(input, ":")

	if len(sp) == 0 || input == "" {
		showPageLive()
		return
	}

	if strings.TrimSpace(sp[0]) == "rid" {
		pageLivePlaySelected(strings.TrimSpace(sp[1]))
	} else {
		r, _ := biliAPI.GetRoomNews("", strings.TrimSpace(sp[1]))
		if r.Code != 0 {
			Logger.Error(r.Message)
			showPageLive()
		} else {
			pageLivePlaySelected2Add = true
			pageLivePlaySelected(r.Data.(map[string]interface{})["roomid"])
		}
	}
}

func pageLiveRecommend() {
	if AccountSelected.SESSDATA == "" || AccountSelected.Uid == 0 {
		Logger.Error("you need to add and select account firstly!")
		showPageLive()
		return
	}
	ret, err := biliAPI.GetLiveUserRecommend(AccountSelected.Uid, AccountSelected.SESSDATA, 1)
	if err != nil {
		Logger.Error(err)
		showPageLive()
	}
	if ret.Code != 0 {
		Logger.Error("Not found your!")
		showPageLive()
	}

	NowLive = make([]*Live, 0, len(ret.Data))
	for _, ll := range ret.Data {
		//log.Printf("%+v\n", ll)
		NowLive = append(NowLive, &Live{
			NikeName: ll.UName,
			Cid:      ll.RoomId,
			Title:    ll.Title,
			State:    1,
		})
	}

	pageLivePlaySelected2Recommend = true
	pageLiveSelect(false)
}

func pageLiveBlock() {
	Logger.Warn("unimplemented method!")
	showPageLive()
}

func pageLiveEdit() {
	Logger.Warn("unimplemented method!")
	showPageLive()
}

func pageLiveDelete() {
	Logger.Warn("unimplemented method!")
	showPageLive()
}
