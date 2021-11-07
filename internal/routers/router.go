package routers

import (
	"net/http"

	"github.com/convee/goblog/internal/handler/admin"
	"github.com/convee/goblog/internal/handler/front"
	"github.com/convee/goblog/internal/routers/middleware"

	"github.com/spf13/viper"
)

func InitRouter() *http.ServeMux {
	mux := &http.ServeMux{}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(viper.GetString("system.root")+"/static"))))
	mux.HandleFunc("/", front.Index)
	mux.HandleFunc("/post", front.PostInfo)
	mux.HandleFunc("/page", front.Page)
	mux.HandleFunc("/tag", front.Tag)

	mux.HandleFunc("/admin/login", admin.Login)
	mux.HandleFunc("/admin/register", admin.Register)
	mux.HandleFunc("/admin/signin", admin.Signin)
	mux.HandleFunc("/admin/signup", admin.Signup)

	mux.Handle("/admin/logout", middleware.AuthWithCookie(http.HandlerFunc(admin.Logout)))
	mux.Handle("/admin", middleware.AuthWithCookie(http.HandlerFunc(admin.PostList)))
	mux.Handle("/admin/post/add", middleware.AuthWithCookie(http.HandlerFunc(admin.PostAdd)))
	mux.Handle("/admin/post/save", middleware.AuthWithCookie(http.HandlerFunc(admin.PostSave)))
	mux.Handle("/admin/post/delete", middleware.AuthWithCookie(http.HandlerFunc(admin.PostDelete)))
	mux.Handle("/admin/page", middleware.AuthWithCookie(http.HandlerFunc(admin.PageList)))
	mux.Handle("/admin/page/add", middleware.AuthWithCookie(http.HandlerFunc(admin.PageAdd)))
	mux.Handle("/admin/page/save", middleware.AuthWithCookie(http.HandlerFunc(admin.PageSave)))
	mux.Handle("/admin/page/delete", middleware.AuthWithCookie(http.HandlerFunc(admin.PageDelete)))
	mux.Handle("/admin/category", middleware.AuthWithCookie(http.HandlerFunc(admin.CategoryList)))
	mux.Handle("/admin/category/add", middleware.AuthWithCookie(http.HandlerFunc(admin.CategoryAdd)))
	mux.Handle("/admin/category/save", middleware.AuthWithCookie(http.HandlerFunc(admin.CategorySave)))
	mux.Handle("/admin/category/delete", middleware.AuthWithCookie(http.HandlerFunc(admin.CategoryDelete)))
	mux.Handle("/admin/tag", middleware.AuthWithCookie(http.HandlerFunc(admin.TagList)))

	return mux
}
