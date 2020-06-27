package parser

import (
	"go-crawler/concurrent-crawler/engine"
	"go-crawler/concurrent-crawler/model"
	"io/ioutil"
	"testing"
)

func TestParseProfile(t *testing.T) {
	url := "http://album.zhenai.com/u/108739485"
	contents, err := ioutil.ReadFile("profile_test_data.html")
	if err != nil {
		panic(err)
	}
	result := ParseProfile(contents, url, "wswinny")
	if len(result.Items) != 1 {
		t.Errorf("Items should contain 1 element; but was %v", result.Items)
	}

	actual := result.Items[0]
	expected := engine.Item{
		Url:  "http://album.zhenai.com/u/108739485",
		Type: "zhenai",
		Id:   "108739485",
		Payload: model.Profile{
			Name:       "wswinny",
			Age:        28,
			Height:     159,
			Weight:     49,
			Income:     "5001-8000",
			Gender:     "女",
			Xinzuo:     "处女座",
			Marriage:   "未婚",
			Education:  "大学本科",
			Occupation: "广告/市场",
			Hukou:      "广东梅州",
			House:      "单位宿舍",
			Car:        "未购车",
		},
	}
	if actual != expected {
		t.Errorf("expected %v, but was %v", expected, actual)
	}
}