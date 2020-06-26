package parser

import (
	"concurrent-crawler/engine"
	"regexp"
)

// `` 包裹的字符不会发生转义
const (
	cityRe = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`
    nextPageRe = `href="(http://www.zhenai.com/zhenghun/[^"]+)"`
)

var (
	// 解析出客户的url
	profileRe = regexp.MustCompile(cityRe)
	// 解析出city内容下一页的url，便于继续爬取此city的客户url
	cityUrlRe = regexp.MustCompile(nextPageRe)
)

func ParseCity(contents []byte, _ string) engine.ParseResult {
	matches := profileRe.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}

	for _, m := range matches {
		url := string(m[1])
		// 把string(m[2])拷贝出来，否则由于闭包特性会导致意想不到的结果
		name := string(m[2])
		// result.Items = append(result.Items, "User " + name)  // 不需要存储
		result.Requests = append(
			result.Requests, engine.Request{
				Url: url,
				// ParseProfile解析函数多了个参数，为了不改变其它解析函数的形式，这里利用了函数式编程闭包的概念
				// ParserFunc函数不是此时运行，由于闭包原因，在函数运行时，m的指向早已经发生了变化，
				// 因此需要拷贝string(m[2])的值来确保准确性
				// ParserFunc: func(c []byte) engine.ParseResult {
				// 	// 闭包，函数中返回了函数
				// 	// return ParseProfile(c, string(m[2]))
				// 	return ParseProfile(c, name)
				// },

			    // 重构后
			    ParserFunc: ProfileParser(name),
			})
	}

	// 看到更多的城市信息
	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests,
			engine.Request {
				Url: string(m[1]),
				// ParserFunc:ParseCity,
				ParserFunc: ParseCity,
			})
	}

	return result
}
