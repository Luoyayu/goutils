package biliAPI

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getKey() (flag bool, hash, key []byte) {
	flag = false
	params := map[string]interface{}{
		"appkey": Config.AppKey,
	}
	params["sign"] = signParams(params)
	l := url.Values{}
	c := &http.Client{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}

	r, _ := http.NewRequest("POST",
		strings.Join([]string{Config.API.Login, "api/oauth2/getKey"}, "/")+"?"+l.Encode(), nil)
	if r != nil {
		r.Header.Add("user-agent", "Mozilla/5.0 BiliDroid/5.37.0")
	}

	resp, err := c.Do(r)
	if err != nil {
		log.Fatal(err)
	} else {
		var ret struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
			Data    struct {
				*TokenInfoStruct
			} `json:"data"`
		}
		b, _ := ioutil.ReadAll(resp.Body)
		if Config.Debug {
			log.Println(string(b))
		}
		defer resp.Body.Close()

		_ = json.Unmarshal(b, &ret)

		if ret.Code != 0 {
			log.Fatal("get Key failed! Error message: ", ret.Message)
		} else {
			flag = true
			if Config.Debug {
				log.Println("ret.Hash: ", ret.Data.Hash)
				log.Println("ret.Key: ", ret.Data.Key)
			}
			hash, key = []byte(ret.Data.Hash), []byte(ret.Data.Key)
		}
	}
	return
}

func DoLogOut() {

}

func DoLogin(username, password string) (ret *respStruct) {
	ok, r1, r2 := getKey()
	if ok {
		ret = login(username, rsaEncrypt(password, r1, r2))
		if Config.Debug {
			log.Println(ret)
		}
		if ret.Code == 0 {
			ret.Data.CookieInfo.CookiesMap = map[string]string{}
			for _, ck := range ret.Data.CookieInfo.Cookies {
				//log.Println(ck.Name, "=", ck.Value)
				ret.Data.CookieInfo.CookiesMap[ck.Name] = ck.Value
			}
			return
		}
	}
	return
}

func login(username, password string) (ret *respStruct) {
	params := map[string]interface{}{
		"appkey":   Config.AppKey,
		"password": password,
		"username": username,
	}
	params["sign"] = signParams(params)

	l := url.Values{}
	c := &http.Client{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}

	r, _ := http.NewRequest("POST",
		strings.Join([]string{Config.API.Login, "api/v3/oauth2/login"}, "/")+"?"+l.Encode(), nil)

	if r != nil {
		r.Header.Add("user-agent", "Mozilla/5.0 BiliDroid/5.37.0")
	}
	resp, err := c.Do(r)
	if err != nil {
		log.Fatal(err)
	} else {
		b, _ := ioutil.ReadAll(resp.Body)

		if Config.Debug {
			log.Println("登录返回值")
			log.Println(string(b))
		}
		_ = json.Unmarshal(b, &ret)
	}
	return
}

func rsaEncrypt(password string, hash, key []byte) string {
	block, _ := pem.Decode(key)
	if block == nil {
		log.Fatal("private key error!")
	}
	pubKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	em, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(string(hash)+password))
	if err != nil {
		log.Fatal(err)
	} else {
		if Config.Debug {
			log.Println("RSA+BASE64 加密结果: ")
			log.Println(base64.StdEncoding.EncodeToString(em))
		}
		return base64.StdEncoding.EncodeToString(em)
	}
	return ""
}
