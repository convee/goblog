package front

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/view"
)

func Tag(c *artgo.Context) {
	tags, _ := daos.GetTags()
	data := make(map[string]interface{})
	data["title"] = "标签"
	data["description"] = "柚子吧的博客标签"
	data["tags"] = tags
	view.Render(data, c, "tag")
}
