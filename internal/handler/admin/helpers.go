package admin

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/pkg"
)

func pagination(c *artgo.Context) (page, perPage int) {
	page, _ = strconv.Atoi(c.Query("page"))
	perPage, _ = strconv.Atoi(c.Query("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func adminPageURL(categoryID, tagID string, page int) string {
	query := url.Values{}
	if categoryID != "" {
		query.Set("category_id", categoryID)
	}
	if tagID != "" {
		query.Set("tag_id", tagID)
	}
	query.Set("page", strconv.Itoa(page))
	return "/admin?" + query.Encode()
}

func renderAdminError(c *artgo.Context, status int, message string) {
	backURL := "/admin"
	if status == http.StatusUnauthorized {
		backURL = "/login"
	}
	pkg.AdminRender(map[string]interface{}{"msg": message, "status": status, "back_url": backURL}, c, "401", status)
}

func positiveFormID(c *artgo.Context) (int, bool) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil || id < 1 {
		renderAdminError(c, http.StatusBadRequest, "无效的记录 ID")
		return 0, false
	}
	return id, true
}
