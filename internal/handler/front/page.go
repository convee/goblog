package front

import (
	"net/http"

	"github.com/convee/goblog/internal/pkg/mysql"
	"github.com/convee/goblog/internal/pkg/view"
)

func Page(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	page := mysql.GetPage(id)
	data := make(map[string]interface{})
	data["page"] = page
	view.Render(data, w, "page")
}
