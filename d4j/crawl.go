package d4j

import (
	"fmt"
	"strconv"
	"sync"
)

// CrawAll
func (r *Redis) CrawAll(maxId int64, minId int64, poolCount int64) {
	if maxId == -1 {
		_, maxIdStr, err := GetBookInfoFromRss(true)
		maxId, err = strconv.ParseInt(maxIdStr, 10, 64)
		if err != nil {
			panic(err)
		}
	}
	if poolCount == 0 {
		poolCount = 20
	}

	if minId <= 0 {
		minId = 1
	}

	jobsCnt := maxId
	wg := sync.WaitGroup{}
	var jobsChan = make(chan task, 4)
	task := make([]task, maxId+1)

	var i int64
	for i = 0; i <= poolCount; i++ {
		go func() {
			for j := range jobsChan {
				istr := fmt.Sprint(j.BookId)
				_, book, err := GetBookInfo(istr)
				if err != nil || book == nil {
					fmt.Println(istr, err)
					continue
				} else {
					fmt.Printf("[%5s]: %s-%s%v %v\n", istr, book.Name, book.Author, book.Tags, book.ShareLink.Url+"/#"+book.ShareLink.Key)
					if Config.SanqiuRedis.Enable == true {
						key := "d4j:" + book.ID
						if err = r.Stores(key, book); err != nil {
							panic(err)
						}
					}
				}
				wg.Done()
			}
		}()
	}

	// send task to goroutine
	for i = minId; i <= jobsCnt; i++ {
		task[i].BookId = i
		jobsChan <- task[i]
		wg.Add(1)
	}
	wg.Wait()
}
