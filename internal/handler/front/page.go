package front

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/pkg"
)

func Page(c *artgo.Context) {
	id := c.Query("id")
	page := daos.GetPage(id)
	data := make(map[string]interface{})
	data["title"] = page.Title
	data["description"] = page.Title
	data["page"] = page
	pkg.Render(data, c, "page")
}
