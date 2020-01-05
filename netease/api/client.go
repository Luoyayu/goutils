package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/luoyayu/goutils/enc"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

var (
	magicKey        = []byte("e82ckenh8dichen8")
	magicSeparator  = "-36cd479b6b5-"
	DEBUG           = false
	CookiesFileName = "cookies_gob"

	// FIXED http request params `header`: w/o MUSIC_U in it
	Header       string
	HeaderSigned string
	HeaderMap    map[string]interface{}
)

func debugln(name string, v ...interface{}) {
	if DEBUG {
		fmt.Println("✨", name, ":", v)
	}
}

func debugV(name string, v interface{}) {
	if DEBUG {
		fmt.Printf("✨✨ %s: %+v\n", name, v)
	}
}

var SINGED = false

func init() {
	DEBUG = os.Getenv("DEBUG") == "true"
	cookiesMap, err := LoadCookies(CookiesFileName)
	SINGED = err == nil

	log.Println("log in status:", SINGED)

	if SINGED {
		HeaderMap = map[string]interface{}{"header": map[string]string{
			"MUSIC_U": cookiesMap["MUSIC_U"],
			"os":      "osx",
			"appver":  "2.3.0"}}
		b, _ := json.Marshal(HeaderMap)
		HeaderSigned = string(b)

	} else {
		HeaderMap = map[string]interface{}{"header": map[string]string{
			"os":     "osx",
			"appver": "2.3.0"}}
		b, _ := json.Marshal(HeaderMap)
		Header = string(b)
	}
}

func initCookie() *cookiejar.Jar {
	// channel=netease; os=osx
	jar, _ := cookiejar.New(nil)
	cookies := []*http.Cookie{
		{
			Name:   "channel",
			Value:  "netease",
			Domain: "music.163.com",
			Path:   "/",
		},
		{
			Name:   "os",
			Value:  "osx",
			Domain: "music.163.com",
			Path:   "/",
		},
	}
	u, _ := url.Parse("http://music.163.com")
	jar.SetCookies(u, cookies)
	debugln("init cookies", jar.Cookies(u))
	return jar
}

// EncData :加/解密参数数据
type EncData struct {
	path          string // api with prefix /api/
	dict          string // params dict for en/de*crypt
	decryptedText string // decrypted Post body
	encryptedText string // encrypted post body
}

// NetEaseClient :
type NetEaseClient struct {
	c   *http.Client
	q   *http.Request
	p   *http.Response
	enc *EncData
}

// NewDefaultHttpClient :新建默认客户端
func (r *NetEaseClient) NewDefaultHttpClient() {
	if r == nil {
		panic("NewDefaultHttpClient called by nil(forget new a NetEaseClient?)")
	}
	r.c = &http.Client{
		Jar: initCookie(),
	}
}

// post :仅能发送加密的body
func (r *NetEaseClient) post(url string) (b []byte, err error) {
	if r.enc.encryptedText == "" {
		err = errors.New("no encrypted body data")
		return
	}
	if r.q, err = http.NewRequest("POST", url, bytes.NewBuffer([]byte(r.enc.encryptedText))); err == nil {
		r.q.Header.Set("host", "music.163.com")
		r.q.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/605.1.15 (KHTML, like Gecko)")
		r.q.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		if r.p, err = r.c.Do(r.q); err == nil {
			debugln("Response Header", r.p.Header)

			b, err = ioutil.ReadAll(r.p.Body)
			debugln("Response Body", string(b))

			// call SaveCookies to update user's MUSIC_U
			auser := User{}
			if err = json.Unmarshal(b, &auser); err == nil {
				cookiesMap := map[string]string{}
				if auser.Code == 200 && auser.Token != "" {
					cookiesMap["MUSIC_U"] = auser.Token
					log.Printf("user: %+v\n", auser)
					err = SaveCookies(CookiesFileName, cookiesMap)
				}
			}
		}
	}
	err = errors.Wrap(err, "*NetEaseClient.post")
	return
}

// DoPost :加密r.path和r.dict构成的r.decryptedText, 发送POST请求
func (r *NetEaseClient) DoPost() (ret []byte, err error) {
	var encryptDataHexStr string
	if encryptDataHexStr, err = r.Encrypt(); err == nil {
		r.enc.encryptedText = "params=" + encryptDataHexStr
		debugln("encrypted Data Hex String with params", r.enc.encryptedText)
		if ret, err = r.post("http://music.163.com/e" + r.enc.path[1:]); err == nil {
			debugln("POST Response Body", string(ret))
			retCode := struct {
				Code int64  `json:"code"`
				Msg  string `json:"msg"`
			}{}
			err = json.Unmarshal(ret, &retCode)
			if retCode.Code != 200 {
				return nil, errors.New(retCode.Msg)
			}
		}
	}
	err = errors.Wrap(err, "*NetEaseClient.DoPost")
	return
}

func (r *NetEaseClient) nobodyUseMd5forEncrypt() string {
	return "nobody" + r.enc.path + "use" + r.enc.dict + "md5forencrypt"
}

func (r *NetEaseClient) NewEnc(path string, params map[string]interface{}) (err error) {
	r.enc = &EncData{}
	r.enc.path = path
	if paramsStr, err := json.Marshal(params); err == nil {
		if SINGED {
			r.enc.dict = strings.TrimSuffix(HeaderSigned, "}") + "," + strings.TrimPrefix(string(paramsStr), "{")
		} else {
			debugln("NewEnc Header", Header)
			r.enc.dict = strings.TrimSuffix(Header, "}") + "," + strings.TrimPrefix(string(paramsStr), "{")
		}
		magic := r.nobodyUseMd5forEncrypt()
		md5HexSign := enc.MD5HexStr(magic)
		r.enc.decryptedText = strings.Join([]string{r.enc.path, r.enc.dict, md5HexSign}, magicSeparator)
		if DEBUG {
			debugln("enc.dict", r.enc.dict)
			debugln("md5HexSign", md5HexSign)
			debugln("decryptedText", r.enc.decryptedText)
		}
	}
	err = errors.Wrap(err, "NewEnc")
	return
}

// Decrypt : Decrypt and apply to r.Path, r.dict
func (r *NetEaseClient) Decrypt() (err error) {
	var decryptBytes []byte
	if decryptBytes, err = DecryptParams(r.enc.decryptedText); err == nil {
		decryptStr := string(decryptBytes)
		decryptList := strings.Split(decryptStr, magicSeparator)
		if len(decryptList) < 3 {
			return errors.New(fmt.Sprint("decrypt failed: split body by magicSeparator turning to", len(decryptList), "part"))
		}
		r.enc.path, r.enc.dict = decryptList[0], decryptList[1]
	}
	return
}

// Encrypt Encrypt r.decryptedText by EcbEncrypt; return Encrypted Hex String
func (r *NetEaseClient) Encrypt() (ret string, err error) {
	ret, err = EncryptParams(r.enc.decryptedText)
	err = errors.Wrap(err, "*NetEaseClient.Encrypt")
	return
}

func EncryptParams(text string) (ret string, err error) {
	if text == "" {
		err = errors.Wrap(errors.New("decrypted Text is None"), "EncryptParams")
		return
	}
	ecb := enc.Enc{}
	ecb.New(magicKey, []byte(text))
	var encryptData []byte
	if encryptData, err = ecb.EcbEncrypt(); err == nil {
		ret = strings.ToUpper(hex.EncodeToString(encryptData))
	}
	return
}

func DecryptParams(text string) (ret []byte, err error) {
	if text == "" {
		err = errors.New("decryptedText is None")
		return
	}
	var textBytes []byte
	if textBytes, err = hex.DecodeString(text); err == nil {
		de := enc.Enc{}
		de.New(magicKey, textBytes)
		if ret, err = de.EcbDecrypt(); err == nil {
			debugln("decrypt params: ", string(ret))
		}
	}
	return
}
