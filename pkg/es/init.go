package es

import (
	"blog/pkg/config"
	"context"
	"github.com/olivere/elastic/v7"
	"log"
)

var esClient *elastic.Client

const blogIndex = "blog"
const mapping = `
{
  "mappings": {
    "properties": {
      "title": {
        "type": "text"
      },
      "content": {
        "type": "keyword"
      }
    }
  }
}`
func Init() {
	client, err := elastic.NewClient(elastic.SetURL(config.EsConfig.Urls...))
	if err != nil {
		log.Fatalln("es new client err ", err)
	}
	esClient = client

	createBlogIndex()

}

func createBlogIndex()  {
	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 检测下索引是否存在
	exists, err := esClient.IndexExists(blogIndex).Do(ctx)
	if err != nil {
		log.Fatalln("es exists handle err ", err)
	}
	if !exists {
		// 索引不存在，则创建一个
		_, err := esClient.CreateIndex(blogIndex).BodyString(mapping).Do(ctx)
		if err != nil {
			log.Fatalln("es create index handle err ", err)
		}
	}
}
