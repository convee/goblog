package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/internal/pkg"
	"github.com/convee/goblog/internal/service"
	"github.com/convee/goblog/pkg/logger"
	"go.uber.org/zap"
)

var postService = service.NewPostService()

func PostList(c *artgo.Context) {
	categoryID := c.Query("category_id")
	tagID := c.Query("tag_id")
	page, perPage := pagination(c)
	posts, err := postService.GetPosts(daos.PostParams{
		CategoryId: categoryID,
		TagId:      tagID,
		PerPage:    perPage + 1,
		Page:       page,
	})
	if err != nil {
		renderAdminError(c, http.StatusInternalServerError, "文章列表加载失败，请重试")
		return
	}
	hasNext := len(posts) > perPage
	if hasNext {
		posts = posts[:perPage]
	}
	categories, err := daos.GetCategories()
	if err != nil {
		renderAdminError(c, http.StatusInternalServerError, "分类列表加载失败，请重试")
		return
	}

	data := map[string]interface{}{
		"posts":      posts,
		"categories": categories,
		"page":       page,
		"has_pre":    page > 1,
		"has_next":   hasNext,
		"pre_url":    adminPageURL(categoryID, tagID, page-1),
		"next_url":   adminPageURL(categoryID, tagID, page+1),
	}
	pkg.AdminRender(data, c, "post_list")
}

func PostAdd(c *artgo.Context) {
	data := make(map[string]interface{})
	var post model.Post
	if id := c.Query("id"); id != "" {
		var err error
		post, err = postService.GetPost(id)
		if err != nil {
			renderAdminError(c, http.StatusNotFound, "文章不存在")
			return
		}
	}
	categories, err := daos.GetCategories()
	if err != nil {
		renderAdminError(c, http.StatusInternalServerError, "分类列表加载失败，请重试")
		return
	}
	data["categories"] = categories
	if post.Id > 0 {
		for i := range categories {
			categories[i].Cur = post.CategoryId
		}
		tags, err := postService.GetPostTags(post)
		if err != nil {
			renderAdminError(c, http.StatusInternalServerError, "文章标签加载失败，请重试")
			return
		}
		data["id"] = post.Id
		data["title"] = post.Title
		data["description"] = post.Description
		data["content"] = post.Content
		data["category_id"] = post.CategoryId
		data["tag_ids"] = post.TagIds
		data["tags"] = tags
	}
	pkg.AdminRender(data, c, "post_add")
}

func PostDelete(c *artgo.Context) {
	postID, ok := positiveFormID(c)
	if !ok {
		return
	}
	if _, err := postService.DeletePost(postID); err != nil {
		renderAdminError(c, http.StatusInternalServerError, "文章删除失败，请重试")
		return
	}
	c.Redirect(http.StatusFound, "/admin")
}

func PostSave(c *artgo.Context) {
	var post model.Post
	post.Id, _ = strconv.Atoi(c.PostForm("id"))
	post.Title = strings.TrimSpace(c.PostForm("title"))
	post.Description = strings.TrimSpace(c.PostForm("description"))
	post.Content = c.PostForm("content")
	post.CategoryId, _ = strconv.Atoi(c.PostForm("category"))
	post.Status = 1
	tags := c.PostForm("tags")

	if post.Id < 0 || post.Title == "" || len([]rune(post.Title)) > 200 || len([]rune(post.Description)) > 2000 || post.CategoryId < 1 {
		renderAdminError(c, http.StatusBadRequest, "请填写有效的标题、摘要和分类")
		return
	}
	if _, err := postService.SavePost(post, tags); err != nil {
		logger.Error("save_post", zap.Int("post_id", post.Id), zap.Error(err))
		renderAdminError(c, http.StatusInternalServerError, "文章保存失败，请检查输入后重试")
		return
	}
	c.Redirect(http.StatusFound, "/admin")
}
