package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/convee/goblog/internal/pkg/model"
	"github.com/convee/goblog/internal/pkg/mysql"
	"github.com/convee/goblog/internal/pkg/view"
)

func PageList(w http.ResponseWriter, r *http.Request) {
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if perPage <= 0 {
		perPage = 20
	}
	if page <= 1 {
		page = 1
	}
	pages, err := mysql.GetPages(mysql.PageParams{
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
	view.AdminRender(data, w, "page/list")
}

func PageAdd(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	id := r.FormValue("id")
	var page model.Page
	if len(id) > 0 {
		page = mysql.GetPage(id)
	}
	categories, _ := mysql.GetCategories()
	data["categories"] = categories

	if page.Id > 0 {

		data["id"] = page.Id
		data["title"] = page.Title
		data["content"] = page.Content
	}
	view.AdminRender(data, w, "page/add")
}

func PageDelete(w http.ResponseWriter, r *http.Request) {
	var page model.Page
	page.Id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	_, err := mysql.PageDelete(page)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "删除失败，请重试"
		view.AdminRender(data, w, "401")
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func PageSave(w http.ResponseWriter, r *http.Request) {
	var page model.Page
	page.Id, _ = strconv.Atoi(r.FormValue("id"))
	page.Title = r.FormValue("title")
	page.Content = r.FormValue("content")
	_, err := mysql.PageSave(page)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "添加或修改失败，请重试"
		view.AdminRender(data, w, "401")
		return
	}
	http.Redirect(w, r, "/admin/page", http.StatusFound)
}
