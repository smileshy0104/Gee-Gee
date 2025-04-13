package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine 是一个实现了http.Handler接口的结构体，用于处理HTTP请求。
type Engine struct{}

// ServeHTTP 实现了http.Handler接口，根据请求的URL路径，返回不同的响应。
// *w 用于写入HTTP响应。
// *req 包含了HTTP请求的所有信息。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		// 根路径返回URL.Path的值。
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		// /hello路径返回所有请求头的键值对。
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		// 其他路径返回404错误信息。
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	// 创建一个Engine实例。
	engine := new(Engine)
	// 监听9999端口并使用Engine实例处理请求，如果出错则记录错误日志并退出。
	log.Fatal(http.ListenAndServe(":9999", engine))
	/** TODO 每当有HTTP请求到达时，标准库会调用 ServeHTTP 方法来处理请求。
	ServeHTTP 方法被自动调用是因为 Engine 实现了 http.Handler 接口，
	并且在 http.ListenAndServe 中注册为处理器。每当有HTTP请求到达时，标准库会调用 ServeHTTP 方法来处理请求。
	这种设计使得开发者可以自定义请求的处理逻辑。
	*/
}
