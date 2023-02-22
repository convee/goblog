package admin

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/es"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/internal/pkg"
	"github.com/convee/goblog/pkg/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/convee/goblog/conf"
)

func PostList(c *artgo.Context) {
	categoryId := c.Query("category_id")
	tagId := c.Query("tag_id")
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	page, _ := strconv.Atoi(c.Query("page"))
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
	posts, err := daos.GetPosts(daos.PostParams{
		CategoryId: categoryId,
		TagId:      tagId,
		PerPage:    perPage,
		Page:       page,
	})
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
	pkg.AdminRender(data, c, "post_list")
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

func PostAdd(c *artgo.Context) {
	data := make(map[string]interface{})
	id := c.PostForm("id")
	var post model.Post
	if len(id) > 0 {
		post = daos.GetPost(id)
	}
	categories, _ := daos.GetCategories()
	data["categories"] = categories

	if post.Id > 0 {
		for i := range categories {
			categories[i].Cur = post.CategoryId
		}
		data["id"] = post.Id
		data["title"] = post.Title
		data["description"] = post.Description
		data["content"] = post.Content
		data["category_id"] = post.CategoryId
		data["tag_ids"] = post.TagIds
		data["tags"] = getTags(post)
	}
	pkg.AdminRender(data, c, "post_add")
}

func getTags(post model.Post) string {
	var tags []string
	if len(post.TagIds) > 0 {
		allTags, _ := daos.GetTags()

		tagsById := make(map[int]model.Tag)
		for _, tag := range allTags {
			tagsById[tag.Id] = tag
		}
		for _, tagId := range post.TagIds {
			tags = append(tags, tagsById[tagId].Name)
		}

	}
	return strings.Join(tags, ",")
}

func PostDelete(c *artgo.Context) {
	var post model.Post
	post.Id, _ = strconv.Atoi(c.Req.URL.Query().Get("id"))
	_, err := daos.PostDelete(post)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "删除失败，请重试"
		pkg.AdminRender(data, c, "401")
		return
	}

	if !conf.Conf.Elasticsearch.Disable {
		go es.DeletePost(es.Post{Id: post.Id})

	}
	c.Redirect(http.StatusFound, "/admin")
}

func PostSave(c *artgo.Context) {
	var post model.Post
	post.Id, _ = strconv.Atoi(c.PostForm("id"))
	post.Title = c.PostForm("title")
	post.Description = c.PostForm("description")
	post.Content = c.PostForm("content")
	post.CategoryId, _ = strconv.Atoi(c.PostForm("category"))
	tags := c.PostForm("tags")
	post.TagIds = getTagIds(tags)
	post.Status = 1
	_, err := daos.PostSave(post)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "添加或修改失败，请重试"
		pkg.AdminRender(data, c, "401")
		return
	}
	if !conf.Conf.Elasticsearch.Disable {
		category := daos.GetCategory(post.CategoryId)
		go es.SavePost(es.Post{
			Id:          post.Id,
			Title:       post.Title,
			Description: post.Description,
			Content:     post.Content,
			Tags:        tags,
			Category:    category.Name,
		})
	}
	for _, tagId := range post.TagIds {
		go daos.IncrTagCount(strconv.Itoa(tagId))
	}
	c.Redirect(http.StatusFound, "/admin")
}

func getTagIds(tags string) (tagIds []int) {
	tagNames := strings.Split(tags, ",")
	tagNames = utils.RemoveDuplicateElement(tagNames)
	allTags, _ := daos.GetTags()
	var allTagNames []string
	allTagByName := make(map[string]model.Tag)
	for _, tag := range allTags {
		allTagNames = append(allTagNames, tag.Name)
		allTagByName[tag.Name] = tag
	}
	for _, tagName := range tagNames {
		if utils.StrInArray(tagName, allTagNames) {
			tagIds = append(tagIds, allTagByName[tagName].Id)
		} else {
			var newTag model.Tag
			newTag.Name = tagName
			newTagId, _ := daos.AddTag(newTag)
			if newTagId > 0 {
				tagIds = append(tagIds, newTagId)
			}
		}
	}
	return
}
