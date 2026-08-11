package front

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/pkg"
	"net/http"
)

func Tag(c *artgo.Context) {
	tags, err := daos.GetTags()
	if err != nil {
		c.Error(http.StatusInternalServerError, "标签列表暂时不可用")
		return
	}
	data := make(map[string]interface{})
	data["title"] = "标签"
	data["tags"] = tags
	pkg.Render(data, c, "tag")
}
