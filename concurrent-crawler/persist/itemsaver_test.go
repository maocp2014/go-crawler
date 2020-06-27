package persist

import (
	"context"
	"encoding/json"
	"go-crawler/concurrent-crawler/engine"
	"go-crawler/concurrent-crawler/model"
	"gopkg.in/olivere/elastic.v5"
	"testing"
)

func TestSave(t *testing.T) {
	profile := engine.Item{
		Url:  "http://albnum.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Payload: model.Profile{
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

	client, err := elastic.NewClient(
		elastic.SetSniff(false))

	const index = "dating_test"
	//save expected
	err = Save(client, index, profile)
	if err != nil {
		panic(err)
	}

	//fetch saved item
	resp, err := client.Get().
		Index(index).
		Type(profile.Type).
		Id(profile.Id).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	//%+v会把结构体的字段名也打出来
	t.Logf("%+v", resp)
	//打印出具体内容
	t.Logf("%s", resp.Source)

	//json反序列化成Profile
	var actual engine.Item
	json.Unmarshal(*resp.Source, &actual)
	t.Logf("%+v", actual.Payload)

	actualProfile, _ := model.FromJsonObj(
		actual.Payload)
	actual.Payload = actualProfile
	t.Logf("%+v", actual.Payload)

	if actual != profile {
		t.Errorf("got %v;expected %v", actual, profile)
	}
}