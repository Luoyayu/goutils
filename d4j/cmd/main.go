package main

import _ "github.com/luoyayu/goutils/d4j"

// import "github.com/luoyayu/goutils/d4j"

func main() {
	//d4j.CrawAll(-1)
	/*st := time.Now()
	  jobsCount := 50 // 任务数量
	  group := sync.WaitGroup{}
	  var jobsChan = make(chan int, 5) // 逻辑处理器数量
	  // a) 生成指定数目的 goroutine，每个 goroutine 消费 jobsChan 中的数据
	  poolCount := 20 // 并发量
	  for i := 0; i < poolCount; i++ {
	  	go func(id int) {
	  		for j := range jobsChan {
	  			fmt.Printf("hello %d in pool %d\n", j, id)
	  			time.Sleep(time.Second)
	  			group.Done()
	  		}
	  	}(i)
	  }

	  // b) 把 job 依次推送到 jobsChan 供 goroutine 消费
	  for i := 0; i < jobsCount; i++ {
	  	jobsChan <- i // 推送任务
	  	group.Add(1)
	  	fmt.Printf("index: %d, goroutine Num: %d\n", i, runtime.NumGoroutine()-1)
	  }
	  group.Wait()
	  fmt.Println("done!", time.Now().Sub(st))*/
}
