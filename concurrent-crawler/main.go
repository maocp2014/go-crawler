package main

import (
	"go-crawler/concurrent-crawler/engine"
	"go-crawler/concurrent-crawler/persist"
	"go-crawler/concurrent-crawler/scheduler"
	"go-crawler/concurrent-crawler/zhenai/parser"
)

func main() {
	// 1、简单串行版SimpleEngine
	// engine.SimpleEngine{}.Run(engine.Request{
	// 	Url:        "http://www.zhenai.com/zhenghun",
	// 	ParserFunc: parser.ParseCityList,
	// })


	// 2、并发版
	// ConcurrentEngineV3兼容了simple scheduler和queued scheduler

	itemChan, err := persist.ItemSaver("data_profile")
	if err != nil {
		panic(err)
	}

	e := engine.ConcurrentEngineV3{
		Scheduler: &scheduler.QueuedSchedulerV2{},
		WorkerCount: 100,
		ItemChan: itemChan,
	}

	e.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}
