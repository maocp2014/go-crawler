package model

import "encoding/json"

// 需要提取解析的字段
type Profile struct {
	// url, id是爬虫比较通用的字段，可以存在types.go
	// Url        string
	// Id         string
	Name       string
	Gender     string
	Age        int
	Height     int
	Weight     int
	Income     string
	Marriage   string
	Education  string
	Occupation string
	Hukou      string
	Xinzuo     string
	House      string
	Car        string
}

func FromJsonObj(o interface{}) (Profile, error) {
	var profile Profile
	// 将数据编码成json字符串
	s, err := json.Marshal(o)
	if err != nil {
		return profile, err
	}

	// 将json字符串解码成Profile
	err = json.Unmarshal(s, &profile)
	return profile, err
}