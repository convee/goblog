package pkg

import (
	"fmt"
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/security"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/convee/goblog/conf"
)

var FuncMap = template.FuncMap{
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

var (
	frontTemplates     map[string]*template.Template
	frontTemplatesErr  error
	frontTemplatesOnce sync.Once
	assetVersion       = fmt.Sprintf("%x", time.Now().Unix())
)

func loadFrontTemplates() error {
	frontTemplatesOnce.Do(func() {
		frontTemplates = make(map[string]*template.Template)
		for _, name := range []string{"index", "post", "page", "tag", "404"} {
			tpl, err := template.New("layout.html").Funcs(FuncMap).ParseFiles(
				"templates/default/layout.html",
				"templates/default/"+name+".html",
			)
			if err != nil {
				frontTemplatesErr = fmt.Errorf("parse front template %s: %w", name, err)
				return
			}
			frontTemplates[name] = tpl
		}
	})
	return frontTemplatesErr
}

func Render(data map[string]interface{}, c *artgo.Context, tpl string) {
	if err := loadFrontTemplates(); err != nil {
		log.Println("front template error:", err)
		c.Error(http.StatusInternalServerError, "页面暂时不可用")
		return
	}
	data["name"] = conf.Conf.App.Name
	data["cdn"] = conf.Conf.App.Cdn
	data["keywords"] = conf.Conf.App.Keywords
	data["navigation"] = conf.Conf.App.Navigation
	data["year"] = time.Now().Year()
	data["asset_version"] = assetVersion
	data["show_nav_search"] = tpl != "index"
	// seo title
	if _, ok := data["title"]; !ok {
		data["title"] = "Go Markdown 博客系统"
	}
	// seo description
	if _, ok := data["description"]; !ok {
		data["description"] = conf.Conf.App.Description
	}
	t, ok := frontTemplates[tpl]
	if !ok {
		c.Error(http.StatusInternalServerError, "页面暂时不可用")
		return
	}
	if err := t.ExecuteTemplate(c.Writer, "layout.html", data); err != nil {
		log.Println("execute front template:", err)
	}
}

func AdminRender(data map[string]interface{}, c *artgo.Context, template string, statuses ...int) {
	data["cdn"] = conf.Conf.App.Cdn
	data["name"] = conf.Conf.App.Name
	csrfToken, err := security.EnsureCSRFToken(c.Writer, c.Req)
	if err != nil {
		c.Error(http.StatusInternalServerError, "页面暂时不可用")
		return
	}
	data["csrf_token"] = csrfToken
	status := http.StatusOK
	if len(statuses) > 0 {
		status = statuses[0]
	}
	c.HTML(status, template+".html", data)
}
