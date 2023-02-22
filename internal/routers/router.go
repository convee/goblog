package routers

import (
	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/handler/admin"
	"github.com/convee/goblog/internal/handler/front"
	"github.com/convee/goblog/internal/pkg"
	"github.com/convee/goblog/internal/routers/middleware"
	"net/http"
)

// InitRouter 自定义 http 服务器
func InitRouter() *artgo.Engine {
	r := artgo.New()
	// 日志中间件，记录全局请求响应日志
	r.Use(middleware.Logger())
	// Recover 中间件，捕获 http server 协程panic错误信息
	r.Use(middleware.Recover())

	// 渲染自定义模板
	r.SetFuncMap(pkg.FuncMap)
	// 加载模板到内存
	r.LoadHTMLGlob("templates/admin/*")
	// 静态文件服务
	r.Static("/static/", "static")

	// 前台路由
	r.GET("/", front.Index)
	r.GET("/favicon.ico", func(c *artgo.Context) {
		http.ServeFile(c.Writer, c.Req, "static/favicon.ico")
	})
	r.GET("/post", front.PostInfo)
	r.GET("/page", front.Page)
	r.GET("/tag", front.Tag)

	// 后台登录路由，不需要校验登录态
	r.GET("/login", admin.Login)
	r.GET("/register", admin.Register)
	r.POST("/signin", admin.Signin)
	r.POST("/signup", admin.Signup)

	// 后台管理页面路由，继承AuthWithCookie中间件，校验登录态
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
