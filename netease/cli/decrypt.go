package main

import (
	"fmt"
	"github.com/luoyayu/netease_go/api"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		panic("please input body encrypted by netease client, like `11675D69CF25E055975....`")
	}
	if ret, err := api.DecryptParams(os.Args[1]); err == nil {
		fmt.Println(string(ret))
	} else {
		panic(err)
	}
}
