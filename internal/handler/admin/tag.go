package admin

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/pkg"
)

func TagList(c *artgo.Context) {
	tags, err := daos.GetTags()
	if err != nil {
		fmt.Println("get ags err:", err)
		return
	}
	data := make(map[string]interface{})
	data["tags"] = tags
	pkg.AdminRender(data, c, "tag_list")
}
