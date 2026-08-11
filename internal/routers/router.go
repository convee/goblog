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
	// 安全响应头与请求体上限应早于业务处理执行。
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.LimitRequestBody())
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
	r.GET("/healthz", front.Health)
	r.GET("/readyz", front.Ready)
	r.GET("/favicon.ico", func(c *artgo.Context) {
		http.ServeFile(c.Writer, c.Req, "static/favicon.ico")
	})
	r.GET("/post/:id", front.PostInfo)
	r.GET("/page/:ident", front.Page)
	r.GET("/tag", front.Tag)

	// 后台登录路由，不需要校验登录态
	r.GET("/login", admin.Login)
	r.POST("/signin", admin.Signin)

	// 后台管理页面路由。每个处理器都由认证和 CSRF 包装器保护，
	// 避免依赖 artgo v1.1.0 无法短路的中间件循环。
	adminGroup := r.Group("/admin")
	adminGroup.POST("/logout", middleware.AdminOnly(admin.Logout))
	adminGroup.GET("/", middleware.AdminOnly(admin.PostList))
	adminGroup.GET("/post/add", middleware.AdminOnly(admin.PostAdd))
	adminGroup.POST("/post/save", middleware.AdminOnly(admin.PostSave))
	adminGroup.POST("/post/delete", middleware.AdminOnly(admin.PostDelete))
	adminGroup.GET("/page", middleware.AdminOnly(admin.PageList))
	adminGroup.GET("/page/add", middleware.AdminOnly(admin.PageAdd))
	adminGroup.POST("/page/save", middleware.AdminOnly(admin.PageSave))
	adminGroup.POST("/page/delete", middleware.AdminOnly(admin.PageDelete))
	adminGroup.GET("/category", middleware.AdminOnly(admin.CategoryList))
	adminGroup.GET("/category/add", middleware.AdminOnly(admin.CategoryAdd))
	adminGroup.POST("/category/save", middleware.AdminOnly(admin.CategorySave))
	adminGroup.POST("/category/delete", middleware.AdminOnly(admin.CategoryDelete))
	adminGroup.GET("/tag", middleware.AdminOnly(admin.TagList))

	return r
}
