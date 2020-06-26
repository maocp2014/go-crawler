package parser

import (
	"concurrent-crawler/engine"
	"concurrent-crawler/model"
	"regexp"
	"strconv"
)

var (
	ageRe = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">([\d]+)岁</div>`)
    heightRe = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">(\d+)cm</div>`)
    incomeRe = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">月收入:([^<]+)</div>`)
    weightRe = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">(\d+)kg</div>`)

    // 没有性别字段
    genderRe = regexp.MustCompile(``)
    xinzuoRe = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">([^<]+)</div>`)
    marriageRe = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">([^<]+)</div>`)
    educationRe = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">([^<]+)</div>`)
    occupationRe = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">([^<]+)</div>`)
    hukouRe = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798="">([^<]+)</div>`)
    houseRe = regexp.MustCompile(`<div class="m-btn pink" data-v-bff6f798="">([^<]+)</div>`)
    carRe = regexp.MustCompile(`<div class="m-btn pink" data-v-bff6f798="">([^<]+)</div>`)

    // 页面出现的"猜你喜欢"的客户
	guessRe = regexp.MustCompile(`<a class="exp-user-name"[^>]*href="(http://albnum.zhenai.com/u/[\d]+)">([^<])</a>`)
	// 用户id
	idUrlRe = regexp.MustCompile(`http://albnum.zhenai.com/u/([\d]+)`)
)

// 解析对应url下名字为name的人的相关信息
func ParseProfile(contents []byte, url string, name string) engine.ParseResult {
	profile := model.Profile{}
	// 名字
	profile.Name = name
	// 年龄
	age, err := strconv.Atoi(extractString(contents, ageRe))
	if err == nil {
		profile.Age = age
	}
	//  身高
	height, err := strconv.Atoi(extractString(contents, heightRe))
	if err == nil {
		profile.Height = height
	}
	// 体重
	weight, err := strconv.Atoi(extractString(contents, weightRe))
	if err == nil {
		profile.Weight = weight
	}
	// 收入
	profile.Income = extractString(contents, incomeRe)
	// 性别
	profile.Gender = extractString(contents, genderRe)
	// 是否有车
	profile.Car = extractString(contents, carRe)
	// 教育状况
	profile.Education = extractString(contents, educationRe)
	// 户口
	profile.Hukou = extractString(contents, hukouRe)
	// 是否有房
	profile.House = extractString(contents, houseRe)
	// 婚姻状况
	profile.Marriage = extractString(contents, marriageRe)
	// 职业
	profile.Occupation = extractString(contents, occupationRe)
	// 星座
	profile.Xinzuo = extractString(contents, xinzuoRe)

	result := engine.ParseResult{
		// Items: []interface{}{profile},
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "zhenai",
				Id:      extractString([]byte(url), idUrlRe),
				Payload: profile,
			},
		},
	}

	// 猜你喜欢
	matches := guessRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		// url := string(m[1])
		// name := string(m[2])

		result.Requests = append(result.Requests,
			engine.Request{
				Url:       string(m[1]),
				// 闭包
				// ParserFunc: func(c []byte) engine.ParseResult {
				// 	return ParseProfile(c, url, name)
				// }

				// 重构后
				ParserFunc: ProfileParser(string(m[2])),
			})
	}
	return result
}

// 从contens中，按照正则表达式提取信息并返回
func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

// 闭包，重新包装parser函数
func ProfileParser(name string) engine.ParserFunc {
	return func(c []byte, url string) engine.ParseResult {
		return ParseProfile(c, url, name)
	}
}