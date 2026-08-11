package admin

import (
	"errors"
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/internal/pkg"
	"net/http"
	"strconv"
	"strings"
)

func CategoryList(c *artgo.Context) {
	categories, err := daos.GetCategories()
	if err != nil {
		renderAdminError(c, http.StatusInternalServerError, "分类列表加载失败，请重试")
		return
	}
	data := make(map[string]interface{})
	data["categories"] = categories
	pkg.AdminRender(data, c, "category_list")
}

func CategoryAdd(c *artgo.Context) {
	data := make(map[string]interface{})
	id, _ := strconv.Atoi(c.Query("id"))
	var category model.Category
	if id > 0 {
		var err error
		category, err = daos.GetCategory(id)
		if err != nil {
			renderAdminError(c, http.StatusNotFound, "分类不存在")
			return
		}
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
	var ok bool
	category.Id, ok = positiveFormID(c)
	if !ok {
		return
	}
	_, err := daos.CategoryDelete(category)
	if err != nil {
		if errors.Is(err, daos.ErrCategoryInUse) {
			renderAdminError(c, http.StatusConflict, "该分类仍有关联文章，无法删除")
			return
		}
		renderAdminError(c, http.StatusInternalServerError, "分类删除失败，请重试")
		return
	}
	c.Redirect(http.StatusFound, "/admin/category")
}

func CategorySave(c *artgo.Context) {
	var category model.Category
	category.Id, _ = strconv.Atoi(c.PostForm("id"))
	category.Name = strings.TrimSpace(c.PostForm("name"))
	if category.Id < 0 || category.Name == "" || len([]rune(category.Name)) > 100 {
		renderAdminError(c, http.StatusBadRequest, "请填写 1 到 100 个字符的分类名称")
		return
	}
	_, err := daos.CategorySave(category)
	if err != nil {
		renderAdminError(c, http.StatusInternalServerError, "分类保存失败，请重试")
		return
	}
	c.Redirect(http.StatusFound, "/admin/category")
}
