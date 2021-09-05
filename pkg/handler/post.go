package handler

import (
	"blog/pkg/model"
	"blog/pkg/mysql"
	"blog/pkg/view"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
)


func Index(w http.ResponseWriter, r *http.Request) {
	categoryId := r.URL.Query().Get("category_id")
	tagId := r.URL.Query().Get("tag_id")
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
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
	posts, err := mysql.GetPosts(mysql.PostParams{
		CategoryId: categoryId,
		TagId:      tagId,
		PerPage:    perPage,
		Page:       page,
	})
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
	return viper.GetString("system.host") + "?" + strings.Join(params, "&")
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
	view.Render(data, w, "post")
}
