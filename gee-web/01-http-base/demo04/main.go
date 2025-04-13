package main

import (
	"fmt"
	"gee-web/gee-web/01-http-base/demo04/gee"
	"net/http"
)

func main() {
	// 创建一个新的 Engine 实例，并使用 GET 和 HEAD 方法注册两个路由处理器。
	r := gee.New()
	// 注册一个 GET 请求的路由处理器，处理根路径 "/" 的请求。
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})
	// 注册一个 GET 请求的路由处理器，处理 "/hello" 路径的请求。
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":9999")
}
