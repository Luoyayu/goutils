package main

import (
	"github.com/luoyayu/goutils/enc"
	"log"
)

func gen(){
	key := "7a840o62v39c41b8"
	sign, err := enc.TimeSignGen("bilibili-api-v3", 600, key, true, "")
	if err != nil {
		panic(err)
	} else {
		log.Println(sign)
	}
}
func check() {
	key := "7a840o62v39c41b8"
	sign := "9wD+Pb8vkwz7mS9uVA5UcXQZnpSc6PTCMTtDZzrqIl3U+H/rdtdtlV6tQp0HW1ia"
	if pass, err := enc.TimeSignCheck(sign, key); err == nil {
		log.Print(pass)
	} else {
		panic(err)
	}
}
func main() {
	gen()
}
