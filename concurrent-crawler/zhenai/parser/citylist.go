package parser

import (
	"fmt"
	"go-crawler/concurrent-crawler/engine"
	"regexp"
)

// 正则运用了技巧
const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func ParseCityList(contents []byte, _ string) engine.ParseResult {
	// <a href="http://www.zhenai.com/zhenghun/aba" data-v-5e16505f="">阿坝</a>
	// re := regexp.MustCompile(`<a href="http://www.zhenai.com/zhenghun/[0-9a-z]+"[^>]*>[^<]+</a>`) // 正则运用了技巧
	// matches := re.FindAll(contents, -1)

	// for _, match := range matches{
	// 	fmt.Printf("%s\n", match)
	// }
	// fmt.Printf("matches found: %d\n", len(matches))

	// 提取城市及其url
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1)

	// for _, m := range matches{
	// 	for _, subMatch := range m {
	// 		fmt.Printf("%s ", subMatch)
	// 	}
	// 	fmt.Println()
	// }

	result := engine.ParseResult{}
	// limit := 10
	for _, m := range matches {
		// result.Items = append(result.Items, "City "+string(m[2]))  // 不需要存储
		result.Requests = append(
			result.Requests, engine.Request{
				Url:        string(m[1]),
				ParserFunc: ParseCity,
			})
		// limit--
		// if limit == 0 {
		// 	break
		// }
		fmt.Printf("City: %s, URL: %s\n", m[2], m[1])
	}
	fmt.Printf("matches found: %d\n", len(matches))
	return result
}
