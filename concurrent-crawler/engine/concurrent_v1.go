package engine

import (
	"log"
)

// version 1
// 所有worker共用一个request input channel
type ConcurrentEngineV1 struct {
	Scheduler   SchedulerV1  // scheduler
	WorkerCount int  // 多少个worker
}

// version 1
// 所有worker共用一个request input channel
// SchedulerV1是1个接口
type SchedulerV1 interface {
	Submit(Request)  // 提交request方法
	ConfigureMasterWorkerChan(chan Request)  // 配置指定worker channel
}

// version 1
// 所有worker共用一个request input channel
func (e *ConcurrentEngineV1) Run(seeds ...Request) {
	// request输入channel
	in := make(chan Request)
	// ParseResult输出channel
	out := make(chan ParseResult)
	// 配置指定的worker channel
	e.Scheduler.ConfigureMasterWorkerChan(in)

	// 创建worker，所有的worker抢同一个channel in中的Request
    for i := 0; i < e.WorkerCount; i++ {
		createWorkerV1(in, out)
	}

	// 把seeds放入request队列
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}
	// 计数
	itemCount := 0

	// 接收PareseResult
	for {
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item #%d: %v\n", itemCount, item)
			itemCount++
		}
		// 把request提交到request队列
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

// version 1
// 所有worker共用一个request input channel
func createWorkerV1(in chan Request, out chan ParseResult) {
	// 每个worker都包裹在一个goroutine中
	go func() {
		for {
			// 从in中取request
			request := <-in
			// 调用worker
			result, err := worker(request)
			if err != nil {
				continue
			}
			// 把result放入out
			out <- result
		}
	}()
}