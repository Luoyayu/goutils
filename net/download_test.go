package net

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestPrintLoadProgress(t *testing.T) {
	url := "http://dldir1.qq.com/qqfile/QQforMac/QQ_6.6.0.dmg"
	_file_ := strings.Split(url, "/")
	fileName := _file_[len(_file_)-1] // cut from url

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	disposition := resp.Header.Get("Content-Disposition")
	sizeStr := resp.Header.Get("Content-Length")

	if len(disposition) != 0 {
		fileName = strings.Split(disposition, "=")[1]
	}

	size, _ := strconv.ParseInt(sizeStr, 10, 64)
	log.Printf("准备下载 %s, 预计大小 %.1f %s\n", fileName, float64(size)/1024/1024, "MBytes")

	f, _ := os.Create(fileName)
	done := make(chan int64)

	go PrintLoadProgress(done, fileName, size)
	n, _ := io.Copy(f, resp.Body)
	done <- n
	defer resp.Body.Close()
}
