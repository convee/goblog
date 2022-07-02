package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/convee/goblog/internal/pkg/model"
	"github.com/convee/goblog/internal/pkg/mysql"
	"github.com/convee/goblog/internal/pkg/view"
)

func CategoryList(w http.ResponseWriter, r *http.Request) {
	categories, err := mysql.GetCategories()
	if err != nil {
		fmt.Println("get categories err:", err)
		return
	}
	data := make(map[string]interface{})
	data["categories"] = categories
	view.AdminRender(data, w, "category/list")
}

func CategoryAdd(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	id, _ := strconv.Atoi(r.FormValue("id"))
	var category model.Category
	if id > 0 {
		category = mysql.GetCategory(id)
	}
	categories, _ := mysql.GetCategories()
	data["categories"] = categories

	if category.Id > 0 {

		data["id"] = category.Id
		data["name"] = category.Name
	}
	view.AdminRender(data, w, "category/add")
}

func CategoryDelete(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	category.Id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	_, err := mysql.CategoryDelete(category)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "删除失败，请重试"
		view.AdminRender(data, w, "401")
		return
	}
	http.Redirect(w, r, "/admin/category", http.StatusFound)
}

func CategorySave(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	category.Id, _ = strconv.Atoi(r.FormValue("id"))
	category.Name = r.FormValue("name")
	_, err := mysql.CategorySave(category)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "添加或修改失败，请重试"
		view.AdminRender(data, w, "401")
		return
	}
	http.Redirect(w, r, "/admin/category", http.StatusFound)
}
