package biliAPI

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type GetCidByAidRet struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    *GetCidByAidRetData `json:"data"`
}

type GetCidByAidRetData struct {
	Aid     int64  `json:"aid"`
	Videos  int    `json:"videos"` // part num
	Pic     string `json:"pic"`
	Title   string `json:"title"`
	PubDate int64  `json:"pubdate"`
	Owner   *struct {
		Mid  int64  `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"owner"`
	Stat *struct {
		View     int `json:"view"`
		Danmaku  int `json:"danmaku"`
		Reply    int `json:"reply"`
		Favorite int `json:"favorite"`
		Coin     int `json:"coin"`
		share    int
		nowRank  int
		hisRank  int
		like     int
		dislike  int
	} `json:"stat"`
	Pages []*struct {
		Cid      int64  `json:"cid"`
		Page     int    `json:"page"`
		From     string `json:"from"`
		Part     string `json:"part"` // part name
		Duration int    `json:"duration"`
	} `json:"pages"`
}

// parts
func GetCidByAid(aid interface{}) (rett *GetCidByAidRet, err error) {
	req := &http.Request{}
	req.Header.Add("Host", "api.bilibili.com")

	if ret, err := GetWithReq(Config.API.GetCIdsByAId, map[string]interface{}{
		"aid": aid,
	}, req, &GetCidByAidRet{}); err == nil {
		rett = ret.(*GetCidByAidRet)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type GetPlayUrlRet struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Quality           int      `json:"quality"`
		Format            string   `json:"format"`
		TimeLength        int64    `json:"timelength"`
		AcceptFormat      string   `json:"accept_format"`
		AcceptDescription []string `json:"accept_description"`
		AcceptQuality     []int    `json:"accept_quality"`
		DUrl              []struct {
			Url    string `json:"url"`
			Length int64  `json:"length"`
			Size   int64  `json:"size"`
		} `json:"durl"`
	} `json:"data"`
}

func GetPlayUrl(aid, cid, qn interface{}, SESSDATA string) (rett *GetPlayUrlRet, err error) {
	/*
		116: 高清1080P60 (需要大会员)
		112: 高清1080P+ (hdflv2) (需要大会员)
		80: 高清1080P (flv)
		74: 高清720P60 (需要大会员)
		64: 高清720P (flv720)
		32: 清晰480P (flv480)
		16: 流畅360P (flv360)
		0:  游客模式(32)
	*/

	req, _ := http.NewRequest("GET",
		"https://api.bilibili.com/x/player/playurl"+"?"+l.Encode(), nil)
	if SESSDATA != "" {
		req.Header.Set("Cookie", "SESSDATA="+SESSDATA)
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36")
	req.Header.Set("Host", "api.bilibili.com")

	if ret, err := GetWithReq(Config.API.GetPlayUrl, map[string]interface{}{
		"avid": aid,
		"cid":  cid,
		"qn":   qn,
	}, req, &GetPlayUrlRet{}); err == nil {
		rett = ret.(*GetPlayUrlRet)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 分块并发下载
// coLimit: 并发数

func DownloadVideo(aid interface{}, p int, qn int, SESSDATA string, coLimit int64) {
	cids, _ := GetCidByAid(aid)

	if p <= 0 { // 全部下载
		for i := 0; i < len(cids.Data.Pages); i++ {

		}
	} else if p <= len(cids.Data.Pages) && p >= 1 { // 下载指定 P
		pCid := cids.Data.Pages[p-1].Cid
		playUrl, _ := GetPlayUrl(aid, pCid, qn, SESSDATA)

		c := &http.Client{}
		size := playUrl.Data.DUrl[0].Size
		title := cids.Data.Pages[p-1].Part

		//done := make(chan int64)

		fileName := fmt.Sprint("av", aid, "-", cids.Data.Title, "-", p, title, ".flv")

		log.Printf("准备下载 %s %s, 预计大小 %.1f %s\n", fileName, playUrl.Data.Format, float64(size)/1024/1024, "MBytes")
		f, _ := os.Create(fileName)
		if err := f.Truncate(size); err != nil {
			panic(err)
		} else {
			log.Println("创建了大小为", size, "字节的空文件")
		}
		defer f.Close()

		//go net.PrintLoadProgress(done, fileName, size)

		if coLimit <= 0 {
			coLimit = 20 // 默认 20 协程
		}
		blockSize := size / coLimit // 每个协程下载大小
		lastSize := size % coLimit  // 最后一个额外协程下载大小

		var i int64
		var wg = sync.WaitGroup{}

		for i = 0; i < coLimit; i++ {
			wg.Add(1)

			left := blockSize * i
			right := blockSize * (i + 1)

			if i == coLimit-1 {
				right += lastSize + 1
			}

			go func(min, max int64, i int64) {

				req, _ := http.NewRequest("GET", playUrl.Data.DUrl[0].Url, nil)

				if SESSDATA != "" {
					req.Header.Set("Cookie", "SESSDATA="+SESSDATA)
				}
				req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:56.0) Gecko/20100101 Firefox/56.0")
				req.Header.Set("Referer", "https://api.bilibili.com/x/web-interface/view?aid="+fmt.Sprint(aid, /*"&p=", p*/))
				req.Header.Set("Origin", "https://www.bilibili.com")

				rangeHeader := "bytes=" + fmt.Sprint(min) + "-" + fmt.Sprint(max-1)
				//log.Println("准备下载第", i, "块 ", rangeHeader)
				req.Header.Set("Range", rangeHeader)
				resp, err := c.Do(req)

				if err != nil {
					panic(err)
				}

				b, _ := ioutil.ReadAll(resp.Body)
				_, _ = f.WriteAt(b, min)
				//log.Println(i, "完成写入", n, "字节")
				b = nil

				resp.Body.Close()
				wg.Done()
			}(left, right, i)
		}
		wg.Wait()
		//done <- 1
	}
}
