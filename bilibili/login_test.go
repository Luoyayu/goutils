package biliAPI

import (
	"fmt"
	"log"
	"testing"
)

func TestDoLogin(t *testing.T) {
	Config.Debug = true
	fmt.Printf("%+v\n", DoLogin("", ""))
}

func Test_getKey(t *testing.T) {
	flag, hash, key := getKey()
	log.Printf("%v %v %v\n", flag, string(hash), string(key))
}
