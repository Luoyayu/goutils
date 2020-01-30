package biliAPI

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type BangumiInitialSate struct {
	MediaInfo MediaInfo `json:"mediaInfo"`
	EpList    []*Ep     `json:"epList"`
	SsList    []*Ss     `json:"ssList"`
}

type MediaInfo struct {
	Id           int64  `json:"id"`
	SsId         int64  `json:"ssId"`
	Title        string `json:"title"`
	JpTitle      string `json:"jpTitle"`
	Series       string `json:"series"`
	Alias        string `json:"alias"`
	Evaluate     string `json:"evaluate"`
	SsType       int    `json:"ssType"`
	SsTypeFormat struct {
		Name     string `json:"name"`
		HomeLink string `json:"homeLink"`
	} `json:"ssTypeFormat"`
	Status      int    `json:"status"`
	SquareCover string `json:"squareCover"`
	Cover       string `json:"cover"`
	Pub         struct {
		Time     string `json:"time"`
		TimeShow string `json:"timeShow"`
		IsStart  bool   `json:"isStart"`
		IsFinish bool   `json:"isFinish"`
	} `json:"pub"`
	Rating struct {
		Score float32 `json:"score"`
		Count int     `json:"count"`
	} `json:"rating"`
	NewestEp struct {
		Id   int64  `json:"id"`
		Desc string `json:"desc"` // "连载中, 每周三 22:30更新"
	} `json:"newestEp"`
	PayMent struct {
		Tip     string `json:"tip"`      // "大会员专享观看特权哦~"
		VipProm string `json:"vip_prom"` // "开通大会员抢先看"
	} `json:"payMent"`
	MainSecTitle string `json:"mainSecTitle"`
}

type Ep struct {
	Id          int64  `json:"id"`
	Badge       string `json:"badge"`    // pay prompt
	EpStatus    int    `json:"epStatus"` // 2 for free; 13 for pay
	Aid         int64  `json:"aid"`
	Cid         int64  `json:"cid"`
	Cover       string `json:"cover"`
	TitleFormat string `json:"titleFormat"`
	LongTitle   string `json:"longTitle"`
	HasNext     bool   `json:"hasNext"`
	I           int    `json:"i"`
}

type Ss struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	PgcType string `json:"pgcType"`
	Cover   string `json:"cover"`
	EpCover string `json:"epCover"`
	Desc    string `json:"desc"`
	Badge   string `json:"badge"`
	Views   int    `json:"views"`
	Follows int    `json:"follows"`
}

func GetBangumiInitialSate(epId interface{}) (*BangumiInitialSate, error) {
	resp, err := http.Get(fmt.Sprint("https://www.bilibili.com/bangumi/play/ep", epId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	var bangumi *BangumiInitialSate
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if strings.HasPrefix(s.Text(), "window.__INITIAL_STATE__=") {
			bangumiText := strings.TrimLeft(s.Text(), "window.__INITIAL_STATE__=")
			bangumiText = bangumiText[:strings.Index(bangumiText, ";(function()")]
			bangumi = &BangumiInitialSate{}
			err = json.Unmarshal([]byte(bangumiText), &bangumi)
			if err != nil {
				bangumi = nil
			}
			return
		}
	})
	return bangumi, err
}
