package admin

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/internal/view"
	"net/http"
	"strconv"
)

func PageList(c *artgo.Context) {
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	page, _ := strconv.Atoi(c.Query("page"))
	if perPage <= 0 {
		perPage = 20
	}
	if page <= 1 {
		page = 1
	}
	pages, err := daos.GetPages(daos.PageParams{
		PerPage: perPage,
		Page:    page,
	})
	if err != nil {
		fmt.Println("get pages err:", err)
		return
	}
	data := make(map[string]interface{})
	data["pages"] = pages
	data["page"] = page
	view.AdminRender(data, c, "page/list")
}

func PageAdd(c *artgo.Context) {
	data := make(map[string]interface{})
	id := c.PostForm("id")
	var page model.Page
	if len(id) > 0 {
		page = daos.GetPage(id)
	}
	categories, _ := daos.GetCategories()
	data["categories"] = categories

	if page.Id > 0 {

		data["id"] = page.Id
		data["title"] = page.Title
		data["content"] = page.Content
	}
	view.AdminRender(data, c, "page/add")
}

func PageDelete(c *artgo.Context) {
	var page model.Page
	page.Id, _ = strconv.Atoi(c.Query("id"))
	_, err := daos.PageDelete(page)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "删除失败，请重试"
		view.AdminRender(data, c, "401")
		return
	}
	c.Redirect(http.StatusFound, "/admin")
}

func PageSave(c *artgo.Context) {
	var page model.Page
	page.Id, _ = strconv.Atoi(c.PostForm("id"))
	page.Title = c.PostForm("title")
	page.Content = c.PostForm("content")
	_, err := daos.PageSave(page)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "添加或修改失败，请重试"
		view.AdminRender(data, c, "401")
		return
	}
	c.Redirect(http.StatusFound, "/admin/page")
}
