package engine

import (
	"concurrent-crawler/fetcher"
	"log"
)

// 串行方式爬虫：最初的实现方式版本，经过重构后，这个文件已经不需要

// 启动 engine
func Run(seeds ...Request) {
	// 请求的切片，维持一个队列
	var requests []Request
	// 循环种子队列
	for _, r := range seeds {
		requests = append(requests, r)
	}

	// 爬取 url
	for len(requests) > 0 {
		// 取第一个 url
		r := requests[0]
		// 切片，保留剩余的
		requests = requests[1:]
		log.Printf("Fetching %s", r.Url)
		// 爬取 url, 返回爬取的内容 []byte
		body, err := fetcher.Fetch(r.Url)
		if err != nil {
			log.Printf("Fetcher: error fetching %s: %v", r.Url, err)
			// 发生错误继续爬取下一个url
			continue
		}
        // 对应各自的解析函数，解析 body
		parseResult := r.ParserFunc(body, r.Url)
		// ...切片展开
		// 把爬取结果里的request继续加到request队列
		requests = append(requests, parseResult.Requests...)

		// 打印解析数据
		for _, item := range parseResult.Items {
			log.Printf("Got item %s", item)
		}
	}
}