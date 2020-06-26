package scheduler

import (
	"concurrent-crawler/engine"
)

// scheduler
type QueuedSchedulerV2 struct {
	// engine goroutine与scheduler goroutine之间通信channel
	// 用于接收request的channel
	requestChan chan engine.Request
	// chan engine.Request 相当是单个worker的channel，所有worker channel则构成workerChan
	// scheduler goroutine与worker goroutine之间通信channel
	// workerChan是channel，对外传输的是chan engine.Request
	workerChan  chan chan engine.Request
}

// 提交Request给requestChan
// 由引擎把request递交给调度器中的channel
// 实现接口方法
func (q *QueuedSchedulerV2) Submit(r engine.Request) {
	q.requestChan <- r
}

// 由外边通知调度器worker准备好了，可以负责接收请求request了
// 在concurrent_v3.go中，实际是由WorkerChan()函数返回的channel传给s.workerChan
// 实现接口方法
func (q *QueuedSchedulerV2) WorkerReady(w chan engine.Request) {
	q.workerChan <- w
}

// 由scheduler来决定worker channel的实现方式，这里是每个worker都有自己的channel
// 由调度器来生成request channel
// 实现接口方法
func (q *QueuedSchedulerV2) WorkerChan() chan engine.Request{
	return make(chan engine.Request)
}

// 实现接口方法
// 运行调度器方法
func (q *QueuedSchedulerV2) Run() {
	// 初始化channel
	q.workerChan = make(chan chan engine.Request)
	q.requestChan = make(chan engine.Request)
    // scheduler goroutine
	go func() {
		// 所有请求队列
		var requestQ []engine.Request
		// 所有工作队列
		var workerQ []chan engine.Request
		for {
			// 可用的request和worker channel
			// 当前请求
			var activeRequest engine.Request
			// 当前worker，用于接收activeRequest
			var activeWorker chan engine.Request
			// 分别从总的请求队列、工作队列中取出第一个队列的元素作为当前活跃的request以及worker
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
            // select...case
			select {
			// requestChan有值
			// 从调度器把新来的请求加入请求队列
			// submit()后，s.requestChan有内容，可执行
			case r := <-q.requestChan:
				requestQ = append(requestQ, r)
			// workerChan有值
			// 从调度器把新来的worker加入worker队列
			// WorkerReady()后，s.workerChan有内容，可执行
			case w := <-q.workerChan:
				workerQ = append(workerQ, w)
			// 如果是当前的请求activeRequest给activeWorker，则此时阻塞在createWorker()中的channel得以运行
			case activeWorker <- activeRequest:
				// 更新总的worker/request队列
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}
