package pkg

import (
	"fmt"
	"github.com/convee/artgo"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/convee/goblog/conf"
)

var FuncMap = template.FuncMap{
	"noescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"formatTime": func(t time.Time, layout string) string {
		return t.Format(layout)
	},
	"tagStyle": func(count int) string {
		var size string
		if count >= 11 {
			size = "60px"
		} else if count == 10 {
			size = "55px"
		} else if count == 9 {
			size = "50px"
		} else if count == 8 {
			size = "45px"
		} else if count == 7 {
			size = "40px"
		} else if count == 6 {
			size = "35px"
		} else if count == 5 {
			size = "30px"
		} else if count == 4 {
			size = "25px"
		} else if count == 3 {
			size = "20px"
		} else if count == 2 {
			size = "15px"
		} else {
			size = "12px"
		}
		return fmt.Sprintf("font-size:%s", size)
	},
}

func Render(data map[string]interface{}, c *artgo.Context, tpl string) {
	var tplPaths []string
	tplPaths = append(tplPaths, "templates/default/layout.html")
	tplPaths = append(tplPaths, "templates/default/"+tpl+".html")
	t, err := template.New("layout.html").Funcs(FuncMap).ParseFiles(tplPaths...)
	if err != nil {
		log.Println("posts template err:", err)
		return
	}
	data["name"] = conf.Conf.App.Name
	data["cdn"] = conf.Conf.App.Cdn
	// seo title
	if _, ok := data["title"]; !ok {
		data["title"] = "Go Markdown 博客系统"
	}
	// seo description
	if _, ok := data["description"]; !ok {
		data["description"] = "Go Markdown 博客系统"
	}
	t.Execute(c.Writer, data)
}

func AdminRender(data map[string]interface{}, c *artgo.Context, template string) {
	data["cdn"] = conf.Conf.App.Cdn
	c.HTML(http.StatusOK, template+".html", data)
}
