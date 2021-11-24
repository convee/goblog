package front

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/internal/pkg/model"
	"github.com/convee/goblog/internal/pkg/mysql"
	"github.com/convee/goblog/internal/pkg/view"
)

func Index(w http.ResponseWriter, r *http.Request) {
	categoryId := r.URL.Query().Get("category_id")
	tagId := r.URL.Query().Get("tag_id")
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	keyword := r.URL.Query().Get("keyword")
	if perPage <= 0 {
		perPage = 20
	}
	if page <= 1 {
		page = 1
	}
	prePage := page - 1
	nextPage := page + 1
	if prePage <= 1 {
		prePage = 1
	}
	params := mysql.PostParams{
		CategoryId: categoryId,
		TagId:      tagId,
		PerPage:    perPage,
		Page:       page,
	}
	// 关键词搜索：标题、描述、内容、分类、标签等
	// todo es search
	if len(keyword) > 0 {
		params.Keyword = keyword
		params.Ids = make(map[string][]string)
		postIds, err := mysql.GetPostIdsByContent(keyword)
		if err == nil {
			params.Ids["ids"] = postIds
		}
		categoryIds, err := mysql.GetCategoryIdsByName(keyword)
		if err == nil {
			params.Ids["category_ids"] = categoryIds
		}
		tagIds, err := mysql.GetTagIdsByName(keyword)
		if err == nil {
			params.Ids["tag_ids"] = tagIds
		}
	}

	posts, err := mysql.GetPosts(params)
	if err != nil {
		fmt.Println("get posts err:", err)
		return
	}
	categories, err := mysql.GetCategories()
	if err != nil {
		fmt.Println("get categories err:", err)
		return
	}
	categoryMap := make(map[int]model.Category)
	for _, category := range categories {
		categoryMap[category.Id] = category
	}
	for index, post := range posts {
		posts[index].CategoryName = categoryMap[post.CategoryId].Name
	}
	//allTags, _ := mysql.GetTags()
	data := make(map[string]interface{})
	data["posts"] = posts
	data["categories"] = categories
	data["page"] = page
	data["pre_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(prePage))
	data["next_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(nextPage))
	view.Render(data, w, "index")
}

func getPageUrl(categoryId string, tagId string, page string) string {
	var params []string
	if len(categoryId) > 0 {
		params = append(params, "category_id="+categoryId)
	}
	if len(tagId) > 0 {
		params = append(params, "tag_id="+tagId)
	}
	params = append(params, "page="+page)
	return conf.Conf.App.Host + "?" + strings.Join(params, "&")
}

func PostInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := r.URL.Query().Get("id")
	post := mysql.GetPost(id)
	category := mysql.GetCategory(post.CategoryId)
	allTags, _ := mysql.GetTags()
	tagIds := post.TagIds
	tagsById := make(map[int]model.Tag)
	for _, tag := range allTags {
		tagsById[tag.Id] = tag
	}
	var tags []model.Tag
	for _, tagId := range tagIds {
		tags = append(tags, tagsById[tagId])
	}
	mysql.IncrView(id)
	post.CategoryName = category.Name
	data := make(map[string]interface{})
	data["post"] = post
	data["tags"] = tags
	data["title"] = post.Title
	view.Render(data, w, "post")
}
