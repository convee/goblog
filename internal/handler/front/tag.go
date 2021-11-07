package front

import (
	"net/http"

	"github.com/convee/goblog/internal/pkg/mysql"
	"github.com/convee/goblog/internal/pkg/view"
)

func Tag(w http.ResponseWriter, r *http.Request) {
	tags, _ := mysql.GetTags()
	data := make(map[string]interface{})
	data["tags"] = tags
	view.Render(data, w, "tag")
}
