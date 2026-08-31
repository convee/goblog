package front

import (
	"net/url"
	"strconv"

	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/service"
)

func postListQueryFromRequest(c *artgo.Context) service.PostListQuery {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))

	return service.NewPostListQuery(
		c.Query("category_id"),
		c.Query("tag_id"),
		c.Query("keyword"),
		page,
		perPage,
	)
}

func pageURL(query service.PostListQuery, page int) string {
	query = service.NewPostListQuery(
		query.CategoryID,
		query.TagID,
		query.Keyword,
		page,
		query.PerPage,
	)

	values := url.Values{}
	if query.CategoryID != "" {
		values.Set("category_id", query.CategoryID)
	}
	if query.TagID != "" {
		values.Set("tag_id", query.TagID)
	}
	if query.Keyword != "" {
		values.Set("keyword", query.Keyword)
	}
	values.Set("page", strconv.Itoa(query.Page))
	values.Set("per_page", strconv.Itoa(query.PerPage))
	return "/?" + values.Encode()
}

func clearFilterURL() string {
	return "/"
}
