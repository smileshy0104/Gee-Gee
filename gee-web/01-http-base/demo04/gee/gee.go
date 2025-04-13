package gee

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc 定义了 Gee 框架中使用的请求处理程序。
// 参数：
//   - http.ResponseWriter: 用于向客户端写入响应。
//   - *http.Request: 包含了 HTTP 请求的相关信息。
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 实现了 ServeHTTP 接口，是 Gee 框架的核心结构体。
// 成员变量：
//   - router: 一个 map，用于存储路由和对应的处理函数。
type Engine struct {
	router map[string]HandlerFunc
}

// New 是 Engine 的构造函数，用于创建一个新的 Engine 实例。
// 返回值：
//   - *Engine: 返回一个指向 Engine 结构体的指针。
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// addRoute 是一个内部方法，用于向路由表中添加路由。
// 参数：
//   - method: HTTP 方法（如 GET、POST）。
//   - pattern: 路由路径。
//   - handler: 对应的处理函数。
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	log.Printf("Route %4s - %s", method, pattern) // 记录路由日志
	engine.router[key] = handler
}

// GET 用于定义处理 GET 请求的路由。
// 参数：
//   - pattern: 路由路径。
//   - handler: 对应的处理函数。
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodGet, pattern, handler)
}

// POST 用于定义处理 POST 请求的路由。
// 参数：
//   - pattern: 路由路径。
//   - handler: 对应的处理函数。
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodPost, pattern, handler)
}

// Delete 添加一个处理DELETE请求的路由
// 参数:
//   - pattern: 路由路径。
//   - handler: 对应的处理函数。
func (engine *Engine) Delete(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodDelete, pattern, handler)
}

// Put 添加一个处理PUT请求的路由
// 参数:
//   - pattern: 路由路径。
//   - handler: 对应的处理函数。
func (engine *Engine) Put(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodPut, pattern, handler)
}

// Patch 添加一个处理PATCH请求的路由
// 参数:
//   - pattern: 路由路径。
//   - handler: 对应的处理函数。
func (engine *Engine) Patch(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodPatch, pattern, handler)
}

// Options 添加一个处理OPTIONS请求的路由
// 参数:
//   - pattern: 路由路径。
//   - handler: 对应的处理函数。
func (engine *Engine) Options(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodOptions, pattern, handler)
}

// Head 添加一个处理HEAD请求的路由
// 参数:
//   - pattern: 路由路径。
//   - handler: 对应的处理函数。
func (engine *Engine) Head(pattern string, handler HandlerFunc) {
	engine.addRoute(http.MethodHead, pattern, handler)
}

// Run 用于启动 HTTP 服务器。
// 参数：
//   - addr: 服务器监听的地址和端口。
//
// 返回值：
//   - error: 如果服务器启动失败，则返回错误信息。
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 实现了 http.Handler 接口，用于处理 HTTP 请求。
// 参数：
//   - w: http.ResponseWriter，用于向客户端写入响应。
//   - req: *http.Request，包含 HTTP 请求的相关信息。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok { // 检查是否存在匹配的路由
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL) // 如果没有匹配的路由，返回 404 错误
	}
}
