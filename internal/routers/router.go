package routers

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/handler/admin"
	"github.com/convee/goblog/internal/handler/front"
	"github.com/convee/goblog/internal/routers/middleware"
	"net/http"
)

func InitRouter() *artgo.Engine {
	r := artgo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

	r.Static("/static/", "static")
	r.GET("/", front.Index)
	r.GET("/favicon.ico", func(c *artgo.Context) {
		http.ServeFile(c.Writer, c.Req, "static/favicon.ico")
	})
	r.GET("/post", front.PostInfo)
	r.GET("/page", front.Page)
	r.GET("/tag", front.Tag)

	r.GET("/login", admin.Login)
	r.GET("/register", admin.Register)
	r.POST("/signin", admin.Signin)
	r.POST("/signup", admin.Signup)

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AuthWithCookie())
	adminGroup.POST("/logout", admin.Logout)
	adminGroup.GET("/", admin.PostList)
	adminGroup.GET("/post/add", admin.PostAdd)
	adminGroup.POST("/post/save", admin.PostSave)
	adminGroup.POST("/post/delete", admin.PostDelete)
	adminGroup.GET("/page", admin.PageList)
	adminGroup.GET("/page/add", admin.PageAdd)
	adminGroup.POST("/page/save", admin.PageSave)
	adminGroup.POST("/page/delete", admin.PageDelete)
	adminGroup.GET("/category", admin.CategoryList)
	adminGroup.GET("/category/add", admin.CategoryAdd)
	adminGroup.POST("/category/save", admin.CategorySave)
	adminGroup.POST("/category/delete", admin.CategoryDelete)
	adminGroup.GET("/tag", admin.TagList)

	return r
}
