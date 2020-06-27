package engine

import (
	"fmt"
	"go-crawler/concurrent-crawler/fetcher"
	"log"
	"strings"
)

// 将厡engine部分功能拆分成worker，方便对其进行并发
func worker(r Request) (ParseResult, error) {
	log.Printf("Fetching %s", r.Url)
	if strings.Contains(r.Url, "qishi") {
		return ParseResult{}, fmt.Errorf("parse %s is wrong,so continue", r.Url)
	}
	body, err := fetcher.Fetch(r.Url)

	if err != nil {
		log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
		return ParseResult{}, err
	}
	// 返回解析结果
	return r.ParserFunc(body, r.Url), nil
}
