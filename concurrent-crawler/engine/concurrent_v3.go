package engine

// 该engine版本统一了所有的engine及scheduler版本

// 兼容各个engine及scheduler版本
type ConcurrentEngineV3 struct {
	Scheduler   SchedulerV3   // 调度器
	WorkerCount int   // 开启的worker goroutine个数
	// ItemChan    chan interface{}  // 存储item通信的channel
	ItemChan    chan Item  // 用于与elasticSearch通信的channel
}

// worker channel准备完毕并通知的接口
type ReadyNotifier interface{
	WorkerReady(chan Request)  // 用于由外边通知调度器的worker，请求request准备完毕
}

// 统一了各个engine及scheduler版本，增加了通用性
// 调度器需要实现的接口
type SchedulerV3 interface {
	ReadyNotifier  // 接口组合方式，一个准备完毕并通知的接口
	Submit(Request)  // 递交请求给调度器
	// ConfigureMasterWorkerChan(chan Request)
	// 为了兼容统一Scheduler接口，去掉了ConfigureMasterWorkerChan方法
	// 加上了WorkerChan方法，由scheduler决定worker使用哪个channel
	// scheduler决定worker共用1个channel还是每个worker都有1个channel
	WorkerChan() chan Request   // 调度器来生成worker channel
	// WorkerReady(chan Request)  // 单独提出来形成1个接口
	Run()   // 运行调度器
}

// engine的run方法
func (e *ConcurrentEngineV3) Run(seeds ...Request) {
	// 用于传递解析请求的结果的channel
	out := make(chan ParseResult)
	// 运行调度器，根据实现方式生成channel和goroutine
	e.Scheduler.Run()

	// 创建worker协程负责解析in中的请求，并返回解析结果给out
	for i := 0; i < e.WorkerCount; i++ {
		createWorkerV3(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	// 循环遍历传入的request，并递交给调度器，createWorker()阻塞，等待相应请求得以运行
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	// itemCount := 0
	for {
		// out的结果由createWorker()解析后返回
		result := <-out
		// 从结果result中取出item，传递给e.ItemChan
		// 阻塞在ItemSaver()中的channel得以运行
		for _, item := range result.Items {
			// log.Printf("Got item #%d: %v\n", itemCount, item)
			// itemCount++
			// 开goroutine，将item送入ItemChan
			// 从结果result中取出item,传递给e.ItemChan
			// 阻塞在ItemSaver()中的channel得以运行
			go func() { e.ItemChan <- item }()
		}

		for _, request := range result.Requests {
			// url去重
			if isDuplicate(request.Url) {
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

// 工作协程，等待请求，并解析返回结果
func createWorkerV3(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			// tell scheduler i'm ready
			ready.WorkerReady(in)
			// 阻塞等待有请求到来，直到queued_v2.go中的activeWorker <- activeRequest
			request := <-in
			// 解析请求，返回结果集
			result, err := worker(request)

			if err != nil {
				continue
			}
			// 阻塞等待结果集result给至out
			out <- result
		}
	}()
}

// 用于url去重
var visitedUrls = make(map[string]bool)

// 用于url去重
func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}