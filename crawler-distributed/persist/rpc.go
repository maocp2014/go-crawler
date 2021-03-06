package persist

import (
	"go-crawler/concurrent-crawler/engine"
	"go-crawler/concurrent-crawler/persist"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

// 提供RPC服务,用于把Item存储到ELasticSearch
type ItemSaverService struct {
	Client *elastic.Client // elasticsearch client
	Index  string          // index(database name)
}

// 用于rpc调用的方法
// 通过elasticSearch client把item保存到elasticSearch的index(database name)中
func (s *ItemSaverService) Save(
	item engine.Item, result *string) error {
	err := persist.Save(s.Client, s.Index, item)
	log.Printf("Item %v saved.", item)
	if err == nil {
		*result = "ok"
	} else {
		log.Printf("Error saving item %v: %v",
			item, err)
	}
	return err
}