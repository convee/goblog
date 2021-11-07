package admin

import (
	"fmt"
	"log"
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
	view.AdminRender(data, w, "admin/category/add")
}

func CategoryDelete(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	category.Id, _ = strconv.Atoi(r.URL.Query().Get("id"))
	_, err := mysql.CategoryDelete(category)
	if err != nil {
		log.Printf("category delete err %v, info:%v", err, category)
	}
	http.Redirect(w, r, "/admin/category/list", http.StatusFound)
}

func CategorySave(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	category.Id, _ = strconv.Atoi(r.FormValue("id"))
	category.Name = r.FormValue("name")
	_, err := mysql.CategorySave(category)
	if err != nil {
		log.Printf("category save err %v, info:%v", err, category)
	}
	http.Redirect(w, r, "/admin/category", http.StatusFound)
}
