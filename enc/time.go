package enc

import (
	"encoding/base64"
	"fmt"
	"github.com/beevik/ntp"
	"strconv"
	"strings"
	"time"
)

// 简单的时间验证签名生成器
func TimeSignGen(serviceName string, limitSecond int64, key string, online bool, ntpServer string) (signBase64 string, err error) {
	var nowT time.Time

	if ntpServer == "" && online == true {
		ntpServer = "ntp.aliyun.com"
		if nowT, err = ntp.Time(ntpServer); err != nil {
			return
		}
	} else if online == false {
		nowT = time.Now()
	}

	outT := nowT.Unix() + limitSecond
	e := &Enc{}
	// base64(过期时间-服务名-base64(过期时间))
	e.New2(key).New2(fmt.Sprint(outT, "-", serviceName, "-", base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(outT)))))
	sign, err := e.EcbEncrypt()
	signBase64 = base64.StdEncoding.EncodeToString(sign)
	return
}

// 时间签名验证 (手动dog)
func TimeSignCheck(sign string, key string) (pass bool, err error) {
	var encText, dencText []byte
	var nowT time.Time
	var out int64
	e := &Enc{}
	if nowT, err = ntp.Time("ntp.aliyun.com"); err == nil {
		if encText, err = base64.StdEncoding.DecodeString(sign); err == nil {
			if dencText, err = e.New2(key).New2(encText).EcbDecrypt(); err == nil {
				if out, err = strconv.ParseInt(strings.Split(string(dencText), "-")[0], 10, 64); err == nil {
					fmt.Println("当前时间: ", nowT.Unix(), "过期时间: ", out)
					pass = nowT.Unix() <= out
				}
			}
		}
	}
	return
}
