package gee

import (
	"net/http"
)

// router是gee框架中的路由管理器
type router struct {
	// handlers是一个映射，用于存储每种HTTP方法-路径组合对应的处理函数
	handlers map[string]HandlerFunc
}

// newRouter创建并返回一个新的router实例
func newRouter() *router {
	// 初始化handlers映射，并返回router实例
	return &router{handlers: make(map[string]HandlerFunc)}
}

// addRoute用于向router中添加路由
// 参数method: HTTP方法，如GET、POST等
// 参数pattern: 路径模式，如"/users/:id"
// 参数handler: 处理函数，当路由匹配时执行
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 构造路由的唯一键，并将处理函数与之关联
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle根据请求的HTTP方法和路径，找到并执行对应的处理函数
// 参数c: 上下文，包含了请求和响应的所有信息
func (r *router) handle(c *Context) {
	// 构造路由的唯一键
	key := c.Method + "-" + c.Path
	// 检查是否存在与当前请求匹配的处理函数
	if handler, ok := r.handlers[key]; ok {
		// 执行处理函数
		handler(c)
	} else {
		// 如果没有找到匹配的处理函数，返回404错误
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
