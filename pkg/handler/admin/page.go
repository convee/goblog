package admin

import (
	"blog/pkg/model"
	"blog/pkg/mysql"
	"blog/pkg/view"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		PerPage:    perPage,
		Page:       page,
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
		log.Printf("page delete err %v, info:%v", err, page)
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
		log.Printf("page save err %v, info:%v", err, page)
	}
	http.Redirect(w, r, "/admin/page", http.StatusFound)
}


