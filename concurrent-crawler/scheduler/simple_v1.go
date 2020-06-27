package scheduler

import "go-crawler/concurrent-crawler/engine"

// 所有worker共用一个input channel输入request
type SimpleSchedulerV1 struct {
	// Scheduler将Request送给worker, 多个worker共用该channel
	workerChan chan engine.Request
}

// 配置指定worker channel
func (s *SimpleSchedulerV1) ConfigureMasterWorkerChan(c chan engine.Request) {
	s.workerChan = c
}

// 把request放入worker channel
func (s *SimpleSchedulerV1) Submit(r engine.Request) {
	// 实现方式1：函数调用
	// 会造成goroutine循环等待，导致goroutine死锁
	// s.workerChan <- r

	// 实现方式2：每个request开启一个goroutine
	// 解决办法: 为每个request开1个goroutine来分发Request，解决循环等待问题
	// 并发分发Request
	go func() { s.workerChan <- r }()
}