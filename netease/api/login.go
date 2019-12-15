package api

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

func (r *NetEaseClient) LoginByCellphone(phone, password string) (err error) {
	if phone == "" || password == "" {
		return errors.New("no phone or password")
	}
	md5Bytes := md5.Sum([]byte(password))
	password = hex.EncodeToString(md5Bytes[:])
	params := map[string]interface{}{
		"phone":       phone,
		"remember":    "true",
		"password":    password,
		"type":        "1",
		"countrycode": "86",
		"e_r":         false,
	}
	if err = r.NewEnc("/api/login/cellphone", params); err == nil {
		_, err = r.DoPost()
	}
	return err
}
