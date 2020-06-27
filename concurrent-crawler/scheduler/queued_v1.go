package scheduler

import (
	"go-crawler/concurrent-crawler/engine"
)

// scheduler
type QueuedSchedulerV1 struct {
	// requestChan用于Engine将request送给Scheduler，从而存储到request队列
	// engine goroutine与scheduler goroutine之间通信channel
	requestChan chan engine.Request
	// workerChan用于Scheduler将Request送给Worker，每个worker有1个单独对应的channel
	// 因此这里的channel是类型嵌套，chan chan类型，把channel组织成1个大的channel
	// chan engine.Request 相当是单个worker的channel，所有worker channel则构成workerChan
	// scheduler goroutine与worker goroutine之间通信channel
	workerChan  chan chan engine.Request
}

// 提交Request给requestChan
// 实现接口方法
func (q *QueuedSchedulerV1) Submit(r engine.Request) {
	q.requestChan <- r
}

// worker channel ready方法
// 实现接口方法
func (q *QueuedSchedulerV1) WorkerReady(w chan engine.Request) {
	// 把已经ready的worker channel放入workerChan，w是一个chan engine.Request
	// 每个worker单独创建的channel
	q.workerChan <- w
}

// scheduler是一个goroutine，调度requests队列和worker队列，并将指定request与worker绑定
// 要改变channel的内容，所以是指针接收者
// 实现接口方法
func (q *QueuedSchedulerV1) Run() {
	// channel初始化
	q.workerChan = make(chan chan engine.Request)
	q.requestChan = make(chan engine.Request)
    // scheduler goroutine
	go func() {
		// request队列
		var requestQ []engine.Request
		// worker队列，chan engine.Request类型
		var workerQ []chan engine.Request
		// 死循环，一直寻找匹配的request以及worker channel
		for {
			// 可用的request
			var activeRequest engine.Request
			// 可用的worker channel
			var activeWorker chan engine.Request
			// 当两队列中都有值时，选择activeWorker、activeRequest
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			// select...case
			select {
			// 当requestChan有值时
			case r := <-q.requestChan:
				// 放入requestQ队列
				requestQ = append(requestQ, r)
			// 当workerChan有值时
			case w := <-q.workerChan:
				// 放入workerQ队列
				workerQ = append(workerQ, w)
			// Request 发给 Worker
			// 把指定的Request放入指定的worker，实现控制
			// requestQ或workerQ为空时，activeWorker或activeRequest为nil，不会被select
			case activeWorker <- activeRequest:
				// 从队列中拿掉Request和Worker，保留剩余的
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}