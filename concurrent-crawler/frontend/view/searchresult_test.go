package view

import (
	"go-crawler/concurrent-crawler/engine"
	"go-crawler/concurrent-crawler/frontend/model"
	common "go-crawler/concurrent-crawler/model"
	"os"
	"testing"
)

func TestSearchResultView_Render(t *testing.T) {
	view := CreateSearchResultView("template.html")
	out, err := os.Create("template.test.html")

	page := model.SearchResult{}
	page.Hits = 123
	item := engine.Item{
		Url:  "http://albnum.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Payload: common.Profile{
			Name:       "安静的雪",
			Age:        34,
			Height:     162,
			Weight:     57,
			Income:     "3001-5000元",
			Gender:     "女",
			Xinzuo:     "牡羊座",
			Marriage:   "离异",
			Education:  "大学本科",
			Occupation: "人事/行政",
			Hukou:      "山东菏泽",
			House:      "已购房",
			Car:        "未购车",
		},
	}

	for i := 0; i < 10; i++ {
		page.Items = append(page.Items, item)
	}

	err = view.Render(out, page)
	if err != nil {
		panic(err)
	}
}