package biliAPI

import (
	"encoding/base64"
	"github.com/luoyayu/goutils/enc"
	"net/url"
)

type Enc struct {
	E *enc.Enc
}

func (r *Enc) New(key string) *Enc {
	r.E = &enc.Enc{}
	r.E.New([]byte(key), nil)
	return r
}

func (r *Enc) EncParams(text string) (string, error) {
	enText, err := r.E.New2(text).EcbEncrypt()
	return url.QueryEscape(base64.StdEncoding.EncodeToString(enText)), err
}

func (r *Enc) DecParams(text string) (string, error) {
	de64, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	deText, err := r.E.New2(de64).EcbDecrypt()
	return string(deText), err
}

