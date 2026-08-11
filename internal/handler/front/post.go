package front

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/pkg"
	"github.com/convee/goblog/internal/service"
)

var frontService = service.NewFrontService()

func Index(c *artgo.Context) {
	categoryID := c.Query("category_id")
	tagID := c.Query("tag_id")
	keyword := strings.TrimSpace(c.Query("keyword"))
	if len([]rune(keyword)) > 100 {
		keyword = string([]rune(keyword)[:100])
	}
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 50 {
		perPage = 50
	}

	posts, err := frontService.GetHomePosts(categoryID, tagID, keyword, page, perPage+1)
	if err != nil {
		c.Error(http.StatusInternalServerError, "文章列表暂时不可用")
		return
	}
	hasNext := len(posts) > perPage
	if hasNext {
		posts = posts[:perPage]
	}
	categories, err := daos.GetCategories()
	if err != nil {
		c.Error(http.StatusInternalServerError, "分类列表暂时不可用")
		return
	}

	data := map[string]interface{}{
		"posts":      posts,
		"categories": categories,
		"keyword":    keyword,
		"page":       page,
		"has_pre":    page > 1,
		"has_next":   hasNext,
		"pre_url":    frontService.GetPageURL(categoryID, tagID, keyword, page-1),
		"next_url":   frontService.GetPageURL(categoryID, tagID, keyword, page+1),
	}
	pkg.Render(data, c, "index")
}

func PostInfo(c *artgo.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil || postID < 1 {
		c.Error(http.StatusNotFound, "文章不存在")
		return
	}
	post, tags, err := frontService.GetPostDetail(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Error(http.StatusNotFound, "文章不存在")
			return
		}
		c.Error(http.StatusInternalServerError, "文章暂时不可用")
		return
	}
	contentHTML, err := pkg.RenderMarkdown(post.Content)
	if err != nil {
		c.Error(http.StatusInternalServerError, "文章暂时不可用")
		return
	}
	data := map[string]interface{}{
		"post":         post,
		"tags":         tags,
		"content_html": contentHTML,
		"title":        post.Title,
		"description":  post.Description,
	}
	pkg.Render(data, c, "post")
}
