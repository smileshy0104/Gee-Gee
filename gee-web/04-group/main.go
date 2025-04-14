package main

import (
	"gee-web/gee-web/04-group/gee"
	"net/http"
)

// main 是程序的入口函数，用于初始化路由和启动 HTTP 服务器。
// 功能描述：
// - 创建一个新的 Gee 路由实例。
// - 定义多个路由规则，包括 GET 和 POST 请求的处理逻辑。
// - 启动 HTTP 服务器并监听指定端口。
func main() {
	// 创建一个新的 Gee 路由实例。
	r := gee.New()

	// 定义根路径下的路由规则。
	r.GET("/index", func(c *gee.Context) {
		// 响应 HTML 内容，显示 Index 页面。
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	// 创建 /v1 路径组。
	v1 := r.Group("/v1")
	{
		// 定义 /v1 路径下的根路由规则。
		v1.GET("/", func(c *gee.Context) {
			// 响应 HTML 内容，显示 Hello Gee。
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		// 定义 /v1/hello 路由规则，支持查询参数 name。
		v1.GET("/hello", func(c *gee.Context) {
			// 根据查询参数 name 构造响应字符串。
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	// 创建 /v2 路径组。
	v2 := r.Group("/v2")
	{
		// 定义 /v2/hello/:name 路由规则，支持路径参数 name。
		v2.GET("/hello/:name", func(c *gee.Context) {
			// 根据路径参数 name 构造响应字符串。
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		// 定义 /v2/login 路由规则，处理 POST 请求。
		v2.POST("/login", func(c *gee.Context) {
			// 将表单数据以 JSON 格式返回。
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	// 启动 HTTP 服务器，监听端口 9999。
	r.Run(":9999")
}
