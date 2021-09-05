package view

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

var funcMap = template.FuncMap{
	"noescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"formatTime": func(t time.Time, layout string) string {
		return t.Format(layout)
	},
}

func Render(data map[string]interface{}, w http.ResponseWriter, tpl string) {
	var tplPaths []string
	tplPaths = append(tplPaths, viper.GetString("system.root")+"/tpl/default/layout.html")
	tplPaths = append(tplPaths, viper.GetString("system.root")+"/tpl/default/"+tpl+".html")
	t, err := template.New("layout.html").ParseFiles(tplPaths...)
	if err != nil {
		log.Println("posts template err:", err)
		return
	}
	data["cdn"] = viper.GetString("system.cdn")
	t.Execute(w, data)
}

func AdminRender(data map[string]interface{}, w http.ResponseWriter, tpl string) {
	var tplPaths []string
	tplPaths = append(tplPaths, viper.GetString("system.tpl")+"/"+tpl+".html")
	t, err := template.ParseFiles(tplPaths...)
	if err != nil {
		log.Println("posts template err:", err)
		return
	}
	data["cdn"] = viper.GetString("system.cdn")
	t.Execute(w, data)
}
