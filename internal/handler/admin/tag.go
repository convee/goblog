package admin

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/view"
)

func TagList(c *artgo.Context) {
	tags, err := daos.GetTags()
	if err != nil {
		fmt.Println("get ags err:", err)
		return
	}
	data := make(map[string]interface{})
	data["tags"] = tags
	view.AdminRender(data, c, "tag/list")
}
