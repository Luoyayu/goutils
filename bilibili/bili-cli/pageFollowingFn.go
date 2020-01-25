package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	biliAPI "github.com/luoyayu/goutils/bilibili"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

var back2PageFollowing bool

func pageFollowingCmdSelect(uid interface{}) {
	fn := PromptSelect(fmt.Sprint("select for ", uid), []string{
		FollowingCmdSelectVideo,
		FollowingCmdSelectDynamic,
		CMDBack,
		CMDHome,
		CMDExit,
	})
	switch fn {
	case FollowingCmdSelectVideo:
		pageVideoSelect(uid)
	case FollowingCmdSelectDynamic:
		pageFollowingDynamic(uid)
		return
	case CMDBack:
		pageFollowingSelect()
		return
	case CMDExit:
		exitClear()
	}
}

func pageFollowingSelect() {
	followings := pageSyncCmdSyncFollowing(false)
	options := make([]string, 0, len(followings))

	for _, ff := range followings {
		if ff.Fid != 0 {
			options = append(options, fmt.Sprintf("%-9d %s", ff.Fid, ff.NikeName))
		}
	}

	uid := PromptSelect("select following", options, survey.WithPageSize(10))
	//Logger.Info(result)

	if uid == "" || len(strings.Split(uid, " ")) == 0 {
		showPageFollowing()
		return
	}

	pageFollowingCmdSelect(strings.TrimSpace(strings.Split(uid, " ")[0]))
}

func pageVideoSelect(uid interface{}) {
	ret, err := biliAPI.GetSpaceArcSearch(uid, 50, 1, 0)
	if err != nil {
		Logger.Error(err)
		pageVideoSelect(uid)
		return
	}

	if ret.Code != 0 || ret.Data == nil {
		Logger.Error(ret.Message)
		pageVideoSelect(uid)
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

	aid := PromptSelect("select video", options, survey.WithPageSize(10))
	if aid != "" {
		aid = strings.TrimSpace(strings.Split(strings.TrimSpace(fmt.Sprint(aid)), " ")[0])
	} else {
		pageFollowingCmdSelect(uid)
		return
	}
	pageVideoPartSelect(aid, uid)
}

func pageVideoPartSelect(aid, uid interface{}) {
	cIds, err := biliAPI.GetCidByAid(aid)
	if err != nil {
		Logger.Error(err)
		pageVideoPartSelect(aid, uid)
		return
	}

	if cIds.Code != 0 {
		Logger.Error(cIds.Message)
		pageVideoPartSelect(aid, uid)
		return
	}

	partOptions := make([]string, 0, len(cIds.Data.Pages)+3)

	for _, ff := range cIds.Data.Pages {
		partOptions = append(partOptions, fmt.Sprint(ff.Page, ": ", ff.Part, " ", time.Duration(ff.Duration*1000000000)))
	}

	partOptions = append(partOptions, CMDBack, CMDHome, CMDExit)

	part := PromptSelect("select part for "+cIds.Data.Title, partOptions, survey.WithPageSize(len(partOptions)))

	if part == "" {
		pageVideoSelect(uid)
		return
	}

	switch part {
	case CMDBack:
		if back2PageFollowing {
			back2PageFollowing = false
			showPageFollowing()
			return
		} else {
			pageVideoSelect(uid)
			return
		}
	case CMDHome:
		showPageHome()
		return
	case CMDExit:
		exitClear()
	default:
		part = strings.Split(fmt.Sprint(part), ": ")[0]

		if part != "" {
			selectedPartId, _ := strconv.ParseInt(part, 10, 64)
			cid := cIds.Data.Pages[selectedPartId-1].Cid
			pageControlPartVideo(aid, cid, uid)
		}
	}
}

func pageControlPartVideo(aid, cid, uid interface{}) {
	playUrls, err := biliAPI.GetPlayUrl(aid, cid, 0, AccountSelected.SESSDATA) // get all qualities
	if err != nil {
		Logger.Error(err)
		pageControlPartVideo(aid, cid, uid)
		return
	}
	if playUrls.Code != 0 {
		Logger.Error(err)
		pageControlPartVideo(aid, cid, uid)
		return
	}

	controlOptions := make([]string, 0, len(playUrls.Data.AcceptQuality)+3)
	for i := 0; i < len(playUrls.Data.AcceptQuality); i++ {
		controlOptions = append(controlOptions, fmt.Sprint(playUrls.Data.AcceptQuality[i], ": ", playUrls.Data.AcceptDescription[i]))
	}

	controlOptions = append(controlOptions, CMDBack, CMDHome, CMDExit)

	ret := PromptSelect(fmt.Sprint("play control, select quality"), controlOptions, survey.WithPageSize(len(controlOptions)))
	if strings.Contains(ret, ": ") {
		ret = strings.Split(ret, ": ")[0]
		playUrls, _ := biliAPI.GetPlayUrl(aid, cid, ret, AccountSelected.SESSDATA)
		Logger.Info("play: ", playUrls.Data.DUrl[0].Url)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		stopMpvSafely()
		// FIXME: title
		go playVideo(ctx, aid, playUrls.Data.DUrl[0].Url, fmt.Sprintln(cid))

		controlCallback := PromptSelect("control video", []string{
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
			pageVideoPartSelect(aid, uid)
			return
		case CMDHome:
			showPageHome()
			return
		case CMDExit:
			exitClear()
		}
	}
}

func pageFollowingCmdPlay() {
	input := PromptInput(&survey.Input{
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
	if aid == "" {
		showPageFollowing()
		return
	}

	cids, _ := biliAPI.GetCidByAid(aid)
	uid := cids.Data.Owner.Mid
	back2PageFollowing = true
	pageVideoPartSelect(aid, uid)
}

func pageFollowingDynamic(uid interface{}) {
	dynamicItems := make([]string, 0, 10)

	curs := int64(0)
	cnt := 0
	for {
		ret, err := biliAPI.GetDynamicSpaceHistory(uid, curs, false)

		if err != nil {
			panic(err)
		}

		for _, c := range ret.Data.Cards {
			cnt += 1
			ss := ""
			switch c.Desc.Type {
			case 1: // 引用
				ss += fmt.Sprintln("原文: ", c.CardContent.T1.Item.Content) + "  "
				ss += fmt.Sprint("> 引用 ")

				switch c.CardContent.T1.OriginType {
				case 2:
					ss += fmt.Sprintln("动态") + "  " + fmt.Sprintln("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName) + " "
					ss += fmt.Sprintln("> 引用动态内容: ", c.CardContent.T1.OriginContent.T2.Item.Description) + " "
					ss += fmt.Sprintln("> 引用动态发布时间: ", time.Unix(c.CardContent.T1.OriginContent.T2.Item.UploadTime, 0)) + " "
					ss += fmt.Sprintln("> 引用动态附带图片数量: ", c.CardContent.T1.OriginContent.T2.Item.PicturesCount) + " "

					for i, p := range c.CardContent.T1.OriginContent.T2.Item.Pictures {
						ss += fmt.Sprintln("第", i+1, "张: ", p.ImgSrc)
						//fmt.Println("第", i+1, "张: ", p.ImgSrc)
					}
				case 4:
					ss += fmt.Sprintln("图文") + " "
					ss += fmt.Sprintln("> 引用图文内容: ", c.CardContent.T1.OriginContent.T4.Item.Content) + " "
					ss += fmt.Sprintln("图文发布时间: ", time.Unix(c.CardContent.T1.OriginContent.T4.Item.Timestamp, 0)) + " "
				case 8:
					ss += fmt.Sprintln("投稿") + " "
					ss += fmt.Sprintln("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName) + " "
					ss += fmt.Sprintln("> 引用投稿ID: ", c.CardContent.T1.OriginContent.T8.Aid) + " "
					ss += fmt.Sprintln("> 引用投稿标题: ", c.CardContent.T1.OriginContent.T8.Title) + " "
					ss += fmt.Sprintln("> 引用投稿描述: ", c.CardContent.T1.OriginContent.T8.Desc) + " "
					ss += fmt.Sprintln("> 引用投稿封面: ", c.CardContent.T1.OriginContent.T8.Pic) + " "
				case 16:
					ss += fmt.Sprintln("小视频") + " "
					ss += fmt.Sprintln("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName) + " "
					ss += fmt.Sprintln("> 引用小视频ID: ", c.CardContent.T1.OriginContent.T16.Item.Id) + " "
					ss += fmt.Sprintln("> 引用小视频图文: ", c.CardContent.T1.OriginContent.T16.Item.Description) + " "
					ss += fmt.Sprintln("> 引用小视频链接: ", c.CardContent.T1.OriginContent.T16.Item.VideoPlayUrl) + " "
				case 32:
					panic("捕获引用类型32!")
				case 64:
					ss += fmt.Sprintln("专栏") + " "
					ss += fmt.Sprintln("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName) + " "
					ss += fmt.Sprintln("> 引用专栏ID: ", c.CardContent.T1.OriginContent.T64.Id) + " "
					ss += fmt.Sprintln("> 引用专栏动态: ", c.CardContent.T1.OriginContent.T64.Dynamic) + " "
					ss += fmt.Sprintln("> 引用专栏标题: ", c.CardContent.T1.OriginContent.T64.Title) + " "
					ss += fmt.Sprintln("> 引用专栏摘要: ", c.CardContent.T1.OriginContent.T64.Summary) + " "
				case 256:
					ss += fmt.Sprintln("创作") + " "
					ss += fmt.Sprintln("> 引用作品标题: ", c.CardContent.T1.OriginContent.T256.Title) + " "
					ss += fmt.Sprintln("> 引用作品介绍: ", c.CardContent.T1.OriginContent.T256.Intro) + " "
					ss += fmt.Sprintln("> 引用作品类型: ", c.CardContent.T1.OriginContent.T256.TypeInfo) + " "
					ss += fmt.Sprintln("> 引用作者: ", c.CardContent.T1.OriginContent.T256.Author) + " "
				case 512:
					ss += fmt.Sprintln("番剧")
					ss += fmt.Sprintln("> 引用番剧ID: ", c.CardContent.T1.OriginContent.T512.Aid) + " "
					ss += fmt.Sprintln("> 引用番剧名: ", c.CardContent.T1.OriginContent.T512.ApiSeasonInfo.Title) + " "
					ss += fmt.Sprintln("> 引用番剧封面: ", c.CardContent.T1.OriginContent.T512.Cover) + " "
					ss += fmt.Sprintln("> 引用番剧第", c.CardContent.T1.OriginContent.T512.Index, "集") + " "
					ss += fmt.Sprintln("> 引用番剧最新一集描述: ", c.CardContent.T1.OriginContent.T512.NewDesc) + " "
					ss += fmt.Sprintln("> 引用番剧播放量: ", c.CardContent.T1.OriginContent.T512.PlayCount) + " "
				case 1024:
					ss += fmt.Sprintln("已失效资源") + " "
					ss += fmt.Sprintln("> 引用提示: ", c.CardContent.T1.Item.Tips) + " "
				case 2048:
					ss += fmt.Sprintln("宣传") + " "
					ss += fmt.Sprintln("> 引用宣传动态内容", c.CardContent.T1.OriginContent.T2048.Vest.Content) + " "
					ss += fmt.Sprintln("> 引用宣传页标题", c.CardContent.T1.OriginContent.T2048.Sketch.Title) + " "
					ss += fmt.Sprintln("> 引用宣传页封面", c.CardContent.T1.OriginContent.T2048.Sketch.CoverUrl) + " "
					ss += fmt.Sprintln("> 引用宣传页描述", c.CardContent.T1.OriginContent.T2048.Sketch.DescText) + " "
					ss += fmt.Sprintln("> 引用宣传页链接", c.CardContent.T1.OriginContent.T2048.Sketch.TargetUrl) + " "
				case 4099:
					ss += fmt.Sprintln("番剧") + " "
					ss += fmt.Sprintln("> 引用番剧ID: ", c.CardContent.T1.OriginContent.T4099.Aid) + " "
					ss += fmt.Sprintln("> 引用番剧名: ", c.CardContent.T1.OriginContent.T4099.ApiSeasonInfo.Title) + " "
					ss += fmt.Sprintln("> 引用番剧封面: ", c.CardContent.T1.OriginContent.T4099.Cover) + " "
					ss += fmt.Sprintln("> 引用番剧第", c.CardContent.T1.OriginContent.T4099.Index, "集") + " "
					ss += fmt.Sprintln("> 引用番剧最新一集描述: ", c.CardContent.T1.OriginContent.T4099.NewDesc) + " "
					ss += fmt.Sprintln("> 引用番剧播放量: ", c.CardContent.T1.OriginContent.T4099.PlayCount) + " "
				case 4200:
					ss += fmt.Sprintln("房间")
					ss += fmt.Sprintln("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName) + " "
					ss += fmt.Sprintln("> 引用房间ID: ", c.CardContent.T1.OriginContent.T4200.RoomId) + " "
					ss += fmt.Sprintln("> 引用房间标题: ", c.CardContent.T1.OriginContent.T4200.Title) + " "
					ss += fmt.Sprintln("> 引用房间标签: ", c.CardContent.T1.OriginContent.T4200.Tags) + " "
				case 4300:
					ss += fmt.Sprintln("收藏夹") + " "
					ss += fmt.Sprintln("> 引用收藏夹ID : ", c.CardContent.T1.OriginContent.T4300.Fid) + " "
					ss += fmt.Sprintln("> 引用收藏夹标题 : ", c.CardContent.T1.OriginContent.T4300.Title) + " "
					ss += fmt.Sprintln("> 引用收藏夹所属 : ", c.CardContent.T1.OriginContent.T4300.Upper.Name) + " "
				default:
					panic(fmt.Sprint("捕获未知引用类型!", c.CardContent.T1.OriginType, c.CardContent.T1.OriginString))
				}
			case 2: // 图文
				ss += fmt.Sprintln("图文内容: ", c.CardContent.T2.Item.Description) + " "
				ss += fmt.Sprintln("发布时间: ", time.Unix(c.CardContent.T2.Item.UploadTime, 0)) + " "
				ss += fmt.Sprintln("附带图片数量: ", c.CardContent.T2.Item.PicturesCount) + " "
				for i, p := range c.CardContent.T2.Item.Pictures {
					ss += fmt.Sprintln("第", i+1, "张: ", p.ImgSrc) + " "
				}
			case 4:
				ss += fmt.Sprintln("图文内容: ", c.CardContent.T4.Item.Content) + " "
				ss += fmt.Sprintln("图文发布时间: ", time.Unix(c.CardContent.T4.Item.Timestamp, 0)) + " "
			case 8: //投稿
				ss += fmt.Sprintln("投稿ID: ", c.CardContent.T8.Aid) + " "
				ss += fmt.Sprintln("投稿标题: ", c.CardContent.T8.Title) + " "
				ss += fmt.Sprintln("投稿描述: ", c.CardContent.T8.Desc) + " "
				ss += fmt.Sprintln("投稿时间: ", time.Unix(c.CardContent.T8.CTime, 0)) + " "
				//ss += fmt.Sprintln(c.CardContent.JumpUrl)
				ss += fmt.Sprintln("投稿封面: ", c.CardContent.T8.Pic) + " "
			case 16: // 小视频
				ss += fmt.Sprintln("> 小视频ID: ", c.CardContent.T16.Item.Id) + " "
				ss += fmt.Sprintln("> 小视频图文: ", c.CardContent.T16.Item.Description) + " "
				ss += fmt.Sprintln("> 小视频链接: ", c.CardContent.T16.Item.VideoPlayUrl) + " "
			case 32:
				panic("捕获未知动态类型32!")
			case 64: // 专栏
				ss += fmt.Sprintln("专栏ID: ", c.CardContent.T64.Id) + " "
				ss += fmt.Sprintln("专栏标题: ", c.CardContent.T64.Title) + " "
				ss += fmt.Sprintln("专栏摘要: ", c.CardContent.T64.Summary) + " "
				ss += fmt.Sprintln("专栏横幅: ", c.CardContent.T64.BannerUrl) + " "
				ss += fmt.Sprintln("专栏发布日期: ", time.Unix(c.CardContent.T64.PublishTime, 0)) + " "
				ss += fmt.Sprintln("专栏创建日期: ", time.Unix(c.CardContent.T64.CTime, 0)) + " "
				ss += fmt.Sprintln("专栏字数: ", c.CardContent.T64.Words) + " "
				for i, img := range c.CardContent.T64.ImageUrls {
					ss += fmt.Sprintln("第%d张专栏图: %s\n", i+1, img) + " "
				}
			case 256:
				ss += fmt.Sprintln("作品标题: ", c.CardContent.T256.Title) + " "
				ss += fmt.Sprintln("作品介绍: ", c.CardContent.T256.Intro) + " "
				ss += fmt.Sprintln("作品类型: ", c.CardContent.T256.TypeInfo) + " "
				ss += fmt.Sprintln("作者: ", c.CardContent.T256.Author) + " "
			case 512:
				ss += fmt.Sprintln("番剧ID: ", c.CardContent.T512.Aid) + " "
				ss += fmt.Sprintln("番剧名: ", c.CardContent.T512.ApiSeasonInfo.Title) + " "
				ss += fmt.Sprintln("番剧封面: ", c.CardContent.T512.Cover) + " "
				ss += fmt.Sprintln("番剧第", c.CardContent.T512.Index, "集") + " "
				ss += fmt.Sprintln("番剧最新一集描述: ", c.CardContent.T512.NewDesc) + " "
				ss += fmt.Sprintln("番剧播放量: ", c.CardContent.T512.PlayCount) + " "
			case 1024:
				ss += fmt.Sprintln("已失效")
			case 2048:
				ss += fmt.Sprintln("宣传动态内容", c.CardContent.T2048.Vest.Content) + " "
				ss += fmt.Sprintln("宣传页标题", c.CardContent.T2048.Sketch.Title) + " "
				ss += fmt.Sprintln("宣传页封面", c.CardContent.T2048.Sketch.CoverUrl) + " "
				ss += fmt.Sprintln("宣传页描述", c.CardContent.T2048.Sketch.DescText) + " "
				ss += fmt.Sprintln("宣传页链接", c.CardContent.T2048.Sketch.TargetUrl) + " "

			default:
				logrus.Println("动态字符串", c.Desc.Type, c.CardString)
				panic("WAIT!")
			}
			ss += fmt.Sprintf("转发数： %d\t评论数: %d\t点赞数: %d\n", c.Desc.RePost, c.Desc.Comment, c.Desc.Like)
			dynamicItems = append(dynamicItems, ss)
		}
		curs = ret.Data.NextOffset
		if cnt >= 10 {
			break
		}
	}
	if PromptSelect("dynamic for "+fmt.Sprint(uid), dynamicItems, survey.WithPageSize(2)) == "" {
		pageFollowingCmdSelect(uid)
	}

}
