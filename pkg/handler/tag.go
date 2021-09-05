package handler

import (
	"blog/pkg/mysql"
	"blog/pkg/view"
	"net/http"
)

func Tag(w http.ResponseWriter, r *http.Request) {
	tags, _ := mysql.GetTags()
	data := make(map[string]interface{})
	data["tags"] = tags
	view.Render(data, w, "tag")
}
