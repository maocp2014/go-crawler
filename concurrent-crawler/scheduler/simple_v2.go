package scheduler

import "go-crawler/concurrent-crawler/engine"

// 该版本统一scheduler适用于concurrent_v3.go版本

// 所有worker共用一个input channel输入request
type SimpleSchedulerV2 struct {
	// Scheduler将Request送给worker, 多个worker共用该channel
	workerChan chan engine.Request
}

// scheduler确定使用哪个worker channel，这里的实现方式是共用1个channel
// 实现接口方法
func (s *SimpleSchedulerV2) WorkerChan() chan engine.Request {
	return s.workerChan
}

// 该版本scheduler实现该方法为空
// 实现接口方法
func (s *SimpleSchedulerV2) WorkerReady(chan engine.Request) {
	// 不实现方法
}

// 实现接口方法
// 初始化共用的worker channel
func (s *SimpleSchedulerV2) Run() {
	s.workerChan = make(chan engine.Request)
}

// 把request放入worker channel
// 实现接口方法
func (s *SimpleSchedulerV2) Submit(r engine.Request) {
	// 实现方式1：函数调用方式
	// 会造成goroutine循环等待，导致goroutine死锁
	// s.workerChan <- r

	// 实现方式2：goroutine方式
	// 解决办法: 为每个request开1个goroutine来分发Request，解决循环等待问题
	// 并发分发Request
	go func() { s.workerChan <- r }()
}