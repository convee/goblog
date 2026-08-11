package admin

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/internal/pkg"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var validPageIdent = regexp.MustCompile(`^[A-Za-z0-9_-]{1,20}$`)

func PageList(c *artgo.Context) {
	page, perPage := pagination(c)
	pages, err := daos.GetPages(daos.PageParams{
		PerPage: perPage + 1,
		Page:    page,
	})
	if err != nil {
		renderAdminError(c, http.StatusInternalServerError, "页面列表加载失败，请重试")
		return
	}
	hasNext := len(pages) > perPage
	if hasNext {
		pages = pages[:perPage]
	}
	data := make(map[string]interface{})
	data["pages"] = pages
	data["page"] = page
	data["has_pre"] = page > 1
	data["has_next"] = hasNext
	data["pre_url"] = "/admin/page?page=" + strconv.Itoa(page-1)
	data["next_url"] = "/admin/page?page=" + strconv.Itoa(page+1)
	pkg.AdminRender(data, c, "page_list")
}

func PageAdd(c *artgo.Context) {
	data := make(map[string]interface{})
	id := c.Query("id")
	var page model.Page
	if len(id) > 0 {
		var err error
		page, err = daos.GetPage(id)
		if err != nil {
			renderAdminError(c, http.StatusNotFound, "页面不存在")
			return
		}
	}
	if page.Id > 0 {

		data["id"] = page.Id
		data["title"] = page.Title
		data["ident"] = page.Ident
		data["content"] = page.Content
	}
	pkg.AdminRender(data, c, "page_add")
}

func PageDelete(c *artgo.Context) {
	var page model.Page
	var ok bool
	page.Id, ok = positiveFormID(c)
	if !ok {
		return
	}
	_, err := daos.PageDelete(page)
	if err != nil {
		renderAdminError(c, http.StatusInternalServerError, "页面删除失败，请重试")
		return
	}
	c.Redirect(http.StatusFound, "/admin/page")
}

func PageSave(c *artgo.Context) {
	var page model.Page
	page.Id, _ = strconv.Atoi(c.PostForm("id"))
	page.Title = c.PostForm("title")
	page.Ident = c.PostForm("ident")
	page.Content = c.PostForm("content")
	if page.Id < 0 || strings.TrimSpace(page.Title) == "" || !validPageIdent.MatchString(page.Ident) {
		renderAdminError(c, http.StatusBadRequest, "请填写有效的标题和页面标识")
		return
	}
	_, err := daos.PageSave(page)
	if err != nil {
		renderAdminError(c, http.StatusInternalServerError, "页面保存失败，请重试")
		return
	}
	c.Redirect(http.StatusFound, "/admin/page")
}
