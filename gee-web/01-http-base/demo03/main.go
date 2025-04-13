package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine 是一个统一的请求处理器，用于处理所有 HTTP 请求。
type Engine struct{}

// ServeHTTP 实现了 http.Handler 接口，用于根据请求路径分发请求并返回响应。
// 参数：
//   - w: http.ResponseWriter，用于向客户端写入响应数据。
//   - req: *http.Request，包含当前请求的所有信息。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 根据请求路径进行路由匹配，并返回不同的响应内容。
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		// 遍历请求头，将所有的键值对写入响应中。
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		// 如果路径不匹配，返回 404 错误。
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

// main 函数是程序的入口，用于启动 HTTP 服务器。
func main() {
	engine := new(Engine)
	server := http.Server{
		Handler: engine,  // 使用自定义的 Engine 处理器处理所有请求。
		Addr:    ":9999", // 指定服务器监听的地址和端口。
	}

	// 启动 HTTP 服务器，如果发生错误且不是因为服务器关闭导致，则记录错误日志并退出程序。
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
