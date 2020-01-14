package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	biliAPI "github.com/luoyayu/goutils/bilibili"
	"strconv"
	"strings"
	"time"
)

var backToFollowing bool

// return selectedUid
func pageFollowingSelect() (selectedUid interface{}) {
	followings := pageSyncCmdSyncFollowing(false)
	options := make([]string, 0, len(followings))

	for _, ff := range followings {
		if ff.Fid != 0 {
			options = append(options, fmt.Sprintf("%-9d %s", ff.Fid, ff.NikeName))
		}
	}

	result := promptSelect("select following", options, survey.WithPageSize(10))
	//Logger.Info(result)

	if result == "" {
		pageFollowingCmdSelect()
		return
	}

	selectedUid = strings.Split(result, " ")[0]
	return

	//Logger.Info(selectedUid)
}

func pageVideoSelect(uid interface{}) (selectedAid interface{}) {
	if uid == nil {
		pageFollowingCmdSelect()
		return nil
	}

	ret, err := biliAPI.GetSpaceArcSearch(uid, 20, 1, 0)
	if err != nil {
		Logger.Error(err)
		pageFollowingCmdSelect()
	}

	if ret.Code != 0 || ret.Data == nil {
		Logger.Error(ret.Message)
		pageFollowingSelect()
		return
	}

	options := make([]string, 0, len(ret.Data.List.VList))

	for _, v := range ret.Data.List.VList {
		if v.Aid != 0 {
			options = append(
				options,
				fmt.Sprintf("%9d %s %s", v.Aid, time.Unix(v.Created, 0).Format(time.RFC822), v.Title),
			)
		}
	}

	selectedAid = promptSelect("select video", options, survey.WithPageSize(10))
	if selectedAid != "" {
		selectedAid = strings.TrimSpace(strings.Split(strings.TrimSpace(fmt.Sprint(selectedAid)), " ")[0])
	} else {
		aid := pageVideoSelect(uid)
		cid := pageVideoPartSelect(aid, uid)
		pageControlPartVideo(aid, cid, uid)
		return
	}
	Logger.Info("selectedAid: ", selectedAid)
	return
}

func pageVideoPartSelect(aid, uid interface{}) (selectedPartCid interface{}) {
	if cIds, err := biliAPI.GetCidByAid(aid); err != nil {
		Logger.Error(err)
	} else {
		if cIds.Code != 0 {
			Logger.Error(cIds.Message)
		} else {
			partOptions := make([]string, 0, len(cIds.Data.Pages)+3)

			for _, ff := range cIds.Data.Pages {
				partOptions = append(partOptions, fmt.Sprint(ff.Page, ": ", ff.Part, " ", time.Duration(ff.Duration*1000000000)))
			}

			partOptions = append(partOptions, CMDBack, CMDHome, CMDExit)

			selectedPartNumber := promptSelect("select part for "+cIds.Data.Title, partOptions, survey.WithPageSize(len(partOptions)))
			if selectedPartNumber == "" {
				cid := pageVideoPartSelect(aid, uid)
				pageControlPartVideo(aid, cid, uid)
				return
			} else {
				switch selectedPartNumber {
				case CMDBack:
					Logger.Info("backToFollowing:", backToFollowing)
					if backToFollowing {
						backToFollowing = false
						showPageFollowing()
						return
					} else {
						aid := pageVideoSelect(uid)
						cid := pageVideoPartSelect(aid, uid)
						pageControlPartVideo(aid, cid, uid)
						return
					}
					return
				case CMDHome:
					showPageHome()
					return
				case CMDExit:
					exitClear()
				default:
					selectedPartNumber = strings.Split(fmt.Sprint(selectedPartNumber), ": ")[0]
					if selectedPartNumber != "" {
						selectedPartId, _ := strconv.ParseInt(selectedPartNumber, 10, 64)
						selectedPartCid = cIds.Data.Pages[selectedPartId-1].Cid
					} else {
						cid := pageVideoPartSelect(aid, uid)
						pageControlPartVideo(aid, cid, uid)
					}
					return
				}
			}
		}
	}
	return
}

// 视频P 播放控制页面
func pageControlPartVideo(aid, cid, uid interface{}) {

	playUrls, err := biliAPI.GetPlayUrl(aid, cid, 0, AccountSelected.SESSDATA) // get all qualities
	if err != nil {
		Logger.Error(err)
		cid := pageVideoPartSelect(aid, uid)
		pageControlPartVideo(aid, cid, uid)
		return
	}
	if playUrls.Code != 0 {
		Logger.Error(err)
		cid := pageVideoPartSelect(aid, uid)
		pageControlPartVideo(aid, cid, uid)
		return
	}

	controlOptions := make([]string, 0, len(playUrls.Data.AcceptQuality)+3)
	for i := 0; i < len(playUrls.Data.AcceptQuality); i++ {
		controlOptions = append(controlOptions, fmt.Sprint(playUrls.Data.AcceptQuality[i], ": ", playUrls.Data.AcceptDescription[i]))
	}

	controlOptions = append(controlOptions, CMDBack, CMDHome, CMDExit)

	ret := promptSelect(fmt.Sprint("play control, select quality"), controlOptions, survey.WithPageSize(len(controlOptions)))
	if strings.Contains(ret, ": ") {
		ret = strings.Split(ret, ": ")[0]
		playUrls, _ := biliAPI.GetPlayUrl(aid, cid, ret, AccountSelected.SESSDATA)
		Logger.Info("play: ", playUrls.Data.DUrl[0].Url)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		stopMpvSafely()
		// FIXME: title
		go playVideo(ctx, aid, playUrls.Data.DUrl[0].Url, fmt.Sprintln(cid))

		controlCallback := promptSelect("control video", []string{
			CMDBack,
			CMDHome,
			CMDExit,
		})

		switch controlCallback {
		case CMDBack:
			cancel()
			pageControlPartVideo(aid, cid, uid)
			return
		case CMDHome:
			cancel()
			showPageHome()
		case CMDExit:
			exitClear()
		}

	} else {
		switch ret {
		case CMDBack:
			pageControlPartVideo(aid, pageVideoPartSelect(aid, uid), uid) // 回到选P -> self
			return
		case CMDHome:
			showPageHome()
			return
		case CMDExit:
			exitClear()
		}
	}
}
func pageFollowingCmdSelect() {
	uid := pageFollowingSelect()
	aid := pageVideoSelect(uid)
	cid := pageVideoPartSelect(aid, uid)
	pageControlPartVideo(aid, cid, uid)
}

func pageFollowingCmdAdd() {

	input := promptInput(&survey.Input{
		Message: "av:XXXX",
		Default: "",
		Help:    "",
	}, survey.WithValidator(func(ans interface{}) error {
		if strings.HasPrefix(fmt.Sprint(ans), "av:") == false {
			return errors.New("input must start with `av:`")
		}
		return nil
	}))
	aid := strings.TrimSpace(strings.TrimLeft(input, "av:"))
	cids, _ := biliAPI.GetCidByAid(aid)
	uid := cids.Data.Owner.Mid
	backToFollowing = true
	cid := pageVideoPartSelect(aid, uid)
	pageControlPartVideo(aid, cid, uid)
	backToFollowing = false
}
