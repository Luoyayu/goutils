package api

import (
	"log"
	"testing"
)

func TestSaveCookies(t *testing.T) {
	if err := SaveCookies("cookies_gob_test", map[string]string{"MUSIC_U": "XXXXXX",}); err != nil {
		panic(err)
	} else {
		log.Println("save cookie to cookies_gob_test!")
	}
}

func TestLoadCookies(t *testing.T) {
	if ret, err := LoadCookies("cookies_gob_test"); err != nil {
		panic(err)
	} else {
		log.Println("read from cookies_gob_test", ret)
	}
}
