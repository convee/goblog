package admin

import (
	"fmt"
	"net/http"

	"github.com/convee/goblog/internal/pkg/mysql"
	"github.com/convee/goblog/internal/pkg/view"
)

func TagList(w http.ResponseWriter, r *http.Request) {
	tags, err := mysql.GetTags()
	if err != nil {
		fmt.Println("get ags err:", err)
		return
	}
	data := make(map[string]interface{})
	data["tags"] = tags
	view.AdminRender(data, w, "tag/list")
}
