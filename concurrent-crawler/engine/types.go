package engine

// 函数类型
type ParserFunc func(contents []byte, url string) ParseResult

// 请求，包括url以及对应的parse function
type Request struct {
	Url        string
	// ParserFunc func([]byte) ParseResult  // 函数类型，解析函数
	ParserFunc ParserFunc  // url对应的解析函数
}

type SerializedParser struct {
	Name string      // 函数名
	Args interface{} // 参数
}

// 解析返回的内容，包括Requests以及Items(最终输出的内容)
type ParseResult struct {
	Requests []Request
	// Items    []interface{}  // 可存储任意数据类型
	Items    []Item
}

// 爬虫需保存的数据
type Item struct {
	Url  string
	Id   string
	Type string
	Payload  interface{}   // 具体爬取的数据
}

// 处理nil parser情况，方便编译和调试
func NilParser([]byte) ParseResult {
	return ParseResult{}
}
