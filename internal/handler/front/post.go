package front

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/internal/pkg"
	"strconv"
	"strings"

	"github.com/convee/goblog/conf"
)

func Index(c *artgo.Context) {
	categoryId := c.Query("category_id")
	tagId := c.Query("tag_id")
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	page, _ := strconv.Atoi(c.Query("page"))
	keyword := c.Query("keyword")
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
	params := daos.PostParams{
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
		postIds, err := daos.GetPostIdsByContent(keyword)
		if err == nil {
			params.Ids["ids"] = postIds
		}
		categoryIds, err := daos.GetCategoryIdsByName(keyword)
		if err == nil {
			params.Ids["category_ids"] = categoryIds
		}
		tagIds, err := daos.GetTagIdsByName(keyword)
		if err == nil {
			params.Ids["tag_ids"] = tagIds
		}
	}

	posts, err := daos.GetPosts(params)
	if err != nil {
		fmt.Println("get posts err:", err)
		return
	}
	categories, err := daos.GetCategories()
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
	data := make(map[string]interface{})
	data["posts"] = posts
	data["categories"] = categories
	data["page"] = page
	data["pre_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(prePage))
	data["next_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(nextPage))
	pkg.Render(data, c, "index")
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

func PostInfo(c *artgo.Context) {
	id := c.Param("id")
	post := daos.GetPost(id)
	category := daos.GetCategory(post.CategoryId)
	allTags, _ := daos.GetTags()
	tagIds := post.TagIds
	tagsById := make(map[int]model.Tag)
	for _, tag := range allTags {
		tagsById[tag.Id] = tag
	}
	var tags []model.Tag
	for _, tagId := range tagIds {
		tags = append(tags, tagsById[tagId])
	}
	daos.IncrView(id)
	post.CategoryName = category.Name
	data := make(map[string]interface{})
	data["post"] = post
	data["tags"] = tags
	data["title"] = post.Title
	data["description"] = post.Description
	pkg.Render(data, c, "post")
}
