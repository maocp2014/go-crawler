package engine

import (
	"log"
)

// version 2
// 解决Request goroutine 和 Worker goroutine不可控问题，这里引入了worker队列版
// 引入了worker队列，加上之前的request队列，达到将指定的request送到指定的worker，从而实现控制
type ConcurrentEngineV2 struct {
	Scheduler   SchedulerV2  // scheduler
	WorkerCount int    // worker个数
}

// version 2
// 队列版本，每个worker都有1个单独的channel，worker队列
type SchedulerV2 interface {
	Submit(Request)   // 提交request方法
	// 每个worker都有1个channel，不需要配置worker channel
	// ConfigureMasterWorkerChan(chan Request)
	WorkerReady(chan Request)   // worker ready方法
	Run()  // 开启goroutine
}

func (e *ConcurrentEngineV2) Run(seeds ...Request) {
	// 输出结果channel
	out := make(chan ParseResult)
	// 生成channel，等待任务，这里scheduler是1个goroutine
	e.Scheduler.Run()

	// 创建worker
	for i := 0; i < e.WorkerCount; i++ {
		createWorkerV2(out, e.Scheduler)
	}

	// 把seeds提交到request队列
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	// 计数
	itemCount := 0

	// 取结果
	for {
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item #%d: %v\n", itemCount, item)
			itemCount++
		}

		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

// version 2
func createWorkerV2(out chan ParseResult, s SchedulerV2) {
	// 每个worker都有1个自己的channel
	in := make(chan Request)
	// goroutine，包含worker
	go func() {
		for {
			// 告诉Scheduler, worker channel已经创建好了
			s.WorkerReady(in)
			// 取Request
			request := <-in
			// 调用worker
			result, err := worker(request)

			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
