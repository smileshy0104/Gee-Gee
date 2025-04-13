package main

import (
	"fmt"
	"log"
	"net/http"
)

// main 是程序的入口函数。
// 它注册了两个 HTTP 处理函数，并启动了一个监听在 ":9999" 端口的 HTTP 服务器。
// 如果服务器启动失败，log.Fatal 将记录错误并终止程序。
func main() {
	http.HandleFunc("/", indexHandler)           // 注册根路径 "/" 的处理函数。
	http.HandleFunc("/hello", helloHandler)      // 注册 "/hello" 路径的处理函数。
	log.Fatal(http.ListenAndServe(":9999", nil)) // 启动 HTTP 服务器，绑定到 ":9999" 端口。
}

// indexHandler 是一个 HTTP 处理函数，用于响应根路径 "/" 的请求。
// 参数:
//
//	w - http.ResponseWriter，用于向客户端发送响应。
//	req - *http.Request，包含客户端请求的详细信息。
//
// 功能: 将请求的 URL 路径 (req.URL.Path) 以字符串形式返回给客户端。
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// helloHandler 是一个 HTTP 处理函数，用于响应 "/hello" 路径的请求。
// 参数:
//
//	w - http.ResponseWriter，用于向客户端发送响应。
//	req - *http.Request，包含客户端请求的详细信息。
//
// 功能: 遍历请求头 (req.Header)，并将每个键值对以格式化字符串的形式返回给客户端。
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
