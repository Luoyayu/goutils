package biliAPI

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
)

// get params sign
func signParams(params map[string]interface{}) string {
	q := url.Values{}
	for k, v := range params {
		q.Add(k, fmt.Sprint(v)) // a way to sort params by key
	}

	m := md5.New()
	m.Write([]byte(q.Encode() + Config.AppSecret))
	return hex.EncodeToString(m.Sum(nil))
}
