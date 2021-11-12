package admin

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/convee/goblog/conf"
	"github.com/convee/goblog/internal/pkg/es"
	"github.com/convee/goblog/internal/pkg/model"
	"github.com/convee/goblog/internal/pkg/mysql"
	"github.com/convee/goblog/internal/pkg/utils"
	"github.com/convee/goblog/internal/pkg/view"
)

func PostList(w http.ResponseWriter, r *http.Request) {
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
	data := make(map[string]interface{})
	data["posts"] = posts
	data["categories"] = categories
	data["page"] = page
	data["pre_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(prePage))
	data["next_url"] = getPageUrl(categoryId, tagId, strconv.Itoa(nextPage))
	view.AdminRender(data, w, "post/list")
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

func PostAdd(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	id := r.FormValue("id")
	var post model.Post
	if len(id) > 0 {
		post = mysql.GetPost(id)
	}
	categories, _ := mysql.GetCategories()
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
	view.AdminRender(data, w, "post/add")
}

func getTags(post model.Post) string {
	var tags []string
	if len(post.TagIds) > 0 {
		allTags, _ := mysql.GetTags()

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

func PostDelete(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	post.Id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	_, err := mysql.PostDelete(post)
	if err != nil {
		log.Printf("post delete err %v, info:%v", err, post)
		http.Redirect(w, r, "/admin", http.StatusBadRequest)
		return
	}

	if !conf.Conf.Elasticsearch.Disable {
		go es.DeletePost(es.Post{Id: post.Id})

	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func PostSave(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	post.Id, _ = strconv.Atoi(r.FormValue("id"))
	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	post.Content = r.FormValue("content")
	post.CategoryId, _ = strconv.Atoi(r.FormValue("category"))
	tags := r.FormValue("tags")
	post.TagIds = getTagIds(tags)
	post.Status = 1
	_, err := mysql.PostSave(post)
	if err != nil {
		log.Printf("post save err %v, info:%v", err, post)
		http.Redirect(w, r, "/admin", http.StatusBadRequest)
		return
	}
	if !conf.Conf.Elasticsearch.Disable {
		go es.SavePost(es.Post{Id: post.Id, Title: post.Title, Content: post.Content})
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func getTagIds(tags string) (tagIds []int) {
	tagNames := strings.Split(tags, ",")
	tagNames = utils.RemoveDuplicateElement(tagNames)
	allTags, _ := mysql.GetTags()
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
			newTagId, _ := mysql.AddTag(newTag)
			if newTagId > 0 {
				tagIds = append(tagIds, newTagId)
			}
		}
	}
	return
}
