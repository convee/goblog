package front

import (
	"net/http"

	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/pkg"
)

func Tag(c *artgo.Context) {
	defaultFrontHandler.Tag(c)
}

func (h *FrontHandler) Tag(c *artgo.Context) {
	tags, err := h.frontService.GetTags()
	if err != nil {
		c.Error(http.StatusInternalServerError, "标签列表暂时不可用")
		return
	}
	data := make(map[string]interface{})
	data["title"] = "标签"
	data["tags"] = tags
	pkg.Render(data, c, "tag")
}
