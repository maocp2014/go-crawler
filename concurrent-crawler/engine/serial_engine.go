package engine

import (
	"log"
)

// 简单串行版SimpleEngine：将origin_engine.go进行了重构
type SimpleEngine struct{}

// 串行函数调用方式
func (s SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	// 加入队列
	for _, r := range seeds {
		requests = append(requests, r)
	}

	for len(requests) > 0 {
		// 取url
		r := requests[0]
		// 保留剩余的request
		requests = requests[1:]
        // 调用worker
		parseResult, err := worker(r)
		if err != nil {
			continue
		}
		// 将解析结果加入requests队列，这里利用了"..."切片展开
		requests = append(requests, parseResult.Requests...)

		for _, item := range parseResult.Items {
			log.Printf("Got item %s\n", item)
		}
	}
}