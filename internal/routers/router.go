package routers

import (
	"github.com/convee/goblog/internal/handler/admin"
	"github.com/convee/goblog/internal/handler/front"
	"github.com/convee/goblog/internal/routers/middleware"
	"net/http"
)

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func InitRouter() *http.ServeMux {
	mux := &http.ServeMux{}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/", middleware.RecoverWrap(http.HandlerFunc(front.Index)))
	mux.Handle("/favicon.ico", middleware.RecoverWrap(http.HandlerFunc(faviconHandler)))
	mux.Handle("/post", middleware.RecoverWrap(http.HandlerFunc(front.PostInfo)))
	mux.Handle("/page", middleware.RecoverWrap(http.HandlerFunc(front.Page)))
	mux.Handle("/tag", middleware.RecoverWrap(http.HandlerFunc(front.Tag)))

	mux.Handle("/admin/login", middleware.RecoverWrap(http.HandlerFunc(admin.Login)))
	mux.Handle("/admin/register", middleware.RecoverWrap(http.HandlerFunc(admin.Register)))
	mux.Handle("/admin/signin", middleware.RecoverWrap(http.HandlerFunc(admin.Signin)))
	mux.Handle("/admin/signup", middleware.RecoverWrap(http.HandlerFunc(admin.Signup)))

	mux.Handle("/admin/logout", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.Logout))))
	mux.Handle("/admin", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.PostList))))
	mux.Handle("/admin/post/add", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.PostAdd))))
	mux.Handle("/admin/post/save", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.PostSave))))
	mux.Handle("/admin/post/delete", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.PostDelete))))
	mux.Handle("/admin/page", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.PageList))))
	mux.Handle("/admin/page/add", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.PageAdd))))
	mux.Handle("/admin/page/save", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.PageSave))))
	mux.Handle("/admin/page/delete", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.PageDelete))))
	mux.Handle("/admin/category", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.CategoryList))))
	mux.Handle("/admin/category/add", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.CategoryAdd))))
	mux.Handle("/admin/category/save", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.CategorySave))))
	mux.Handle("/admin/category/delete", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.CategoryDelete))))
	mux.Handle("/admin/tag", middleware.RecoverWrap(middleware.AuthWithCookie(http.HandlerFunc(admin.TagList))))

	return mux
}
