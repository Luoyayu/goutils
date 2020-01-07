package net

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func PrintLoadProgress(done <-chan int64, unFinishedFilepath string, totalSize int64) {
	const phase = "#"
	const phaseToT = 50
	const padding = "-"
	const progressFormat = "\rprogress: |%s%s| %v%% [elapsed: %s left: %s %10s]"

	st := time.Now()
	var lastSize, nowSize int64
	var phaseNow int

	for {
		select {
		case <-done:
			if nowSize != totalSize || phaseNow != phaseToT {
				fmt.Printf(progressFormat,
					strings.Repeat(phase, phaseToT),
					strings.Repeat(padding, 0),
					100,
					FormatTime(time.Now().Sub(st)),
					FormatTime(time.Second*0),
					FormatSpeed(int64(float64(totalSize)/time.Now().Sub(st).Seconds())),
				)
				fmt.Println()
			}
			return
		default:
			file, err := os.Open(unFinishedFilepath)
			if err != nil {
				log.Fatalln(err)
			}

			if fi, err := file.Stat(); err != nil {
				log.Fatalln(err)
			} else {
				lastSize = nowSize
				nowSize = fi.Size()
			}

			speed := nowSize - lastSize // Bytes/s
			timeLeft := time.Duration(int64((float64(totalSize-nowSize) / float64(speed)) * 1e9))
			if timeLeft.Hours() < 0 {
				timeLeft = time.Hour * 24
			}
			rate := float64(nowSize) / float64(totalSize)
			phaseNow = int(rate * phaseToT)
			unfinishedPhase := phaseToT - phaseNow

			fmt.Printf(
				progressFormat, strings.Repeat(phase, phaseNow), strings.Repeat("-", unfinishedPhase),
				int(rate*100),
				FormatTime(time.Now().Sub(st)),
				FormatTime(timeLeft),
				FormatSpeed(speed),
			)
		}
		time.Sleep(time.Second)
	}
}

func FormatTime(d time.Duration) string {
	s := (d % time.Minute) / time.Second
	m := (d % time.Hour) / time.Minute
	h := d / time.Hour
	if h == 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func FormatSpeed(speed int64) string { // speed is Bytes/s
	if speed >= (1<<10) && speed < (1<<20) {
		return fmt.Sprintf("%4.1f KB/s", float64(speed)/1024)
	} else if speed >= (1<<20) && speed < (1<<30) {
		return fmt.Sprintf("%4.1f MB/s", float64(speed)/1024/1024)
	} else if speed >= (1 << 30) {
		return fmt.Sprintf("%4.1f GB/s", float64(speed)/1024/1024/1024)
	}
	return fmt.Sprintf("%4.1f  B/s", float64(speed))
}
