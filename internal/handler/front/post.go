package front

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/internal/pkg"
	"github.com/convee/goblog/internal/service"
)

// FrontService keeps HTTP handlers independent of DAO and template details and
// makes page composition testable with a small fake service.
type FrontService interface {
	GetPostList(query service.PostListQuery) (service.PostList, error)
	GetCategories() ([]model.Category, error)
	GetTags() ([]model.Tag, error)
	GetPostDetail(id string) (model.Post, []model.Tag, error)
	GetPage(ident string) (model.Page, error)
}

type FrontHandler struct {
	frontService FrontService
}

func NewFrontHandler(frontService FrontService) *FrontHandler {
	return &FrontHandler{frontService: frontService}
}

var defaultFrontHandler = NewFrontHandler(service.NewFrontService())

func Index(c *artgo.Context) {
	defaultFrontHandler.Index(c)
}

func (h *FrontHandler) Index(c *artgo.Context) {
	query := postListQueryFromRequest(c)
	postList, err := h.frontService.GetPostList(query)
	if err != nil {
		c.Error(http.StatusInternalServerError, "文章列表暂时不可用")
		return
	}
	categories, err := h.frontService.GetCategories()
	if err != nil {
		c.Error(http.StatusInternalServerError, "分类列表暂时不可用")
		return
	}

	pkg.Render(indexViewData(query, postList, categories), c, "index")
}

func indexViewData(query service.PostListQuery, postList service.PostList, categories []model.Category) map[string]interface{} {
	categoryName := ""
	for i := range categories {
		categories[i].Cur = 0
		if strconv.Itoa(categories[i].Id) == query.CategoryID {
			categories[i].Cur = 1
			categoryName = categories[i].Name
		}
	}
	hasContentFilters := query.CategoryID != "" || query.TagID != "" || query.Keyword != ""

	return map[string]interface{}{
		"posts":              postList.Posts,
		"categories":         categories,
		"keyword":            query.Keyword,
		"category_id":        query.CategoryID,
		"category_name":      categoryName,
		"tag_id":             query.TagID,
		"per_page":           postList.PerPage,
		"page":               postList.Page,
		"has_filters":        hasContentFilters,
		"has_active_filters": hasContentFilters || postList.PerPage != 20,
		"has_pre":            postList.HasPrevious,
		"has_next":           postList.HasNext,
		"pre_url":            pageURL(query, postList.Page-1),
		"next_url":           pageURL(query, postList.Page+1),
		"clear_url":          clearFilterURL(),
	}
}

func PostInfo(c *artgo.Context) {
	defaultFrontHandler.PostInfo(c)
}

func (h *FrontHandler) PostInfo(c *artgo.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil || postID < 1 {
		c.Error(http.StatusNotFound, "文章不存在")
		return
	}
	post, tags, err := h.frontService.GetPostDetail(id)
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
