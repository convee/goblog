package handler

import (
	"blog/pkg/mysql"
	"blog/pkg/view"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	page := mysql.GetPage(id)
	data := make(map[string]interface{})
	data["page"] = page
	view.Render(data, w, "page")
}
