package es

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
)

type Post struct {
	Id          int
	Title       string // 标题
	Description string
	Content     string    // 内容
	Created     time.Time // 创建时间
	Category    string
	Tags        string
}

func AddPost(post Post) {
	ctx := context.Background()
	id := strconv.Itoa(post.Id)
	// 使用client创建一个新的文档
	put1, err := esClient.Index().
		Index(blogIndex). // 设置索引名称
		Id(id).           // 设置文档id
		BodyJson(post).   // 指定前面声明的微博内容
		Do(ctx)           // 执行请求，需要传入一个上下文对象
	if err != nil {
		log.Println("es add post err ", err)
		return
	}

	log.Printf("文档Id %s, 索引名 %s\n", put1.Id, put1.Index)
}

func ExistsPost(post Post) bool {
	ctx := context.Background()
	id := strconv.Itoa(post.Id)
	exist, err := esClient.Exists().Index(blogIndex).Id(id).Do(ctx)
	if err != nil {
		log.Println("es exists post err ", err)
		return false
	}
	return exist
}

func GetPost(post Post) {
	ctx := context.Background()
	id := strconv.Itoa(post.Id)
	// 根据id查询文档
	get1, err := esClient.Get().
		Index(blogIndex). // 指定索引名
		Id(id).           // 设置文档id
		Do(ctx)           // 执行请求
	if err != nil {
		log.Println("es get post err ", err)
		return
	}
	if get1.Found {
		log.Printf("文档id=%s 版本号=%d 索引名=%s\n", get1.Id, get1.Version, get1.Index)
	}

	//# 手动将文档内容转换成go struct对象
	msg2 := Post{}
	// 提取文档内容，原始类型是json数据
	data, _ := get1.Source.MarshalJSON()
	// 将json转成struct结果
	json.Unmarshal(data, &msg2)
	// 打印结果
	log.Println(msg2.Title)
}

func UpdatePost(post Post) {
	id := strconv.Itoa(post.Id)
	ctx := context.Background()
	_, err := esClient.Update().
		Index(blogIndex).                           // 设置索引名
		Id(id).                                     // 文档id
		Doc(map[string]interface{}{"retweets": 0}). // 更新retweets=0，支持传入键值结构
		Do(ctx)                                     // 执行ES查询
	if err != nil {
		log.Println("es update post err ", err)
		return
	}
}

func SavePost(post Post) {
	if ExistsPost(post) {
		UpdatePost(post)
	} else {
		AddPost(post)
	}
}

func DeletePost(post Post) {
	id := strconv.Itoa(post.Id)
	ctx := context.Background()
	// 根据id删除一条数据
	_, err := esClient.Delete().
		Index(blogIndex).
		Id(id).
		Do(ctx)
	if err != nil {
		log.Println("es delete post err ", err)
		return
	}
}

func Search(post Post, perPage int, page int) (err error, ids []string) {
	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	// 创建term查询条件，用于精确查询
	termQuery := elastic.NewTermQuery("Content", post.Content)

	searchResult, err := esClient.Search().
		Index(blogIndex).      // 设置索引名
		Query(termQuery).      // 设置查询条件
		Sort("Created", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		//From(page). // 设置分页参数 - 起始偏移量，从第0行记录开始
		//Size(perPage). // 设置分页参数 - 每页大小
		Pretty(true). // 查询结果返回可读性较好的JSON格式
		Do(ctx)       // 执行请求

	if err != nil {
		log.Println("es search post err ", err)
		return
	}

	log.Printf("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits())
	if searchResult.TotalHits() > 0 {
		// 查询结果不为空，则遍历结果
		var b1 Post
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range searchResult.Each(reflect.TypeOf(b1)) {
			// 转换成Post对象
			if t, ok := item.(Post); ok {
				ids = append(ids, strconv.Itoa(t.Id))
			}
		}
	}
	return nil, ids
}
