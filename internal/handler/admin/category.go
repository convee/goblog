package admin

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/internal/pkg"
	"net/http"
	"strconv"
)

func CategoryList(c *artgo.Context) {
	categories, err := daos.GetCategories()
	if err != nil {
		fmt.Println("get categories err:", err)
		return
	}
	data := make(map[string]interface{})
	data["categories"] = categories
	pkg.AdminRender(data, c, "category_list")
}

func CategoryAdd(c *artgo.Context) {
	data := make(map[string]interface{})
	id, _ := strconv.Atoi(c.PostForm("id"))
	var category model.Category
	if id > 0 {
		category = daos.GetCategory(id)
	}
	categories, _ := daos.GetCategories()
	data["categories"] = categories

	if category.Id > 0 {

		data["id"] = category.Id
		data["name"] = category.Name
	}
	pkg.AdminRender(data, c, "category_add")
}

func CategoryDelete(c *artgo.Context) {
	var category model.Category
	category.Id, _ = strconv.Atoi(c.Query("id"))
	_, err := daos.CategoryDelete(category)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "删除失败，请重试"
		pkg.AdminRender(data, c, "401")
		return
	}
	c.Redirect(http.StatusFound, "/admin/category")
}

func CategorySave(c *artgo.Context) {
	var category model.Category
	category.Id, _ = strconv.Atoi(c.PostForm("id"))
	category.Name = c.PostForm("name")
	_, err := daos.CategorySave(category)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "添加或修改失败，请重试"
		pkg.AdminRender(data, c, "401")
		return
	}
	c.Redirect(http.StatusFound, "/admin/category")
}
