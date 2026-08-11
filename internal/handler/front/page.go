package front

import (
	"database/sql"
	"errors"
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/pkg"
	"net/http"
)

func Page(c *artgo.Context) {
	ident := c.Param("ident")

	page, err := daos.GetPageByIdent(ident)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Error(http.StatusNotFound, "页面不存在")
			return
		}
		c.Error(http.StatusInternalServerError, "页面暂时不可用")
		return
	}
	contentHTML, err := pkg.RenderMarkdown(page.Content)
	if err != nil {
		c.Error(http.StatusInternalServerError, "页面暂时不可用")
		return
	}
	data := make(map[string]interface{})
	data["title"] = page.Title
	data["description"] = page.Title
	data["page"] = page
	data["content_html"] = contentHTML
	pkg.Render(data, c, "page")
}
