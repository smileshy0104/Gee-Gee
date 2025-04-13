package gee

import (
	"net/http"
	"strings"
)

// router是gee框架中的路由管理器
type router struct {
	// roots是一个映射，用于存储路由树的根节点
	roots map[string]*node
	// handlers是一个映射，用于存储每种HTTP方法-路径组合对应的处理函数
	handlers map[string]HandlerFunc
}

// newRouter创建并返回一个新的router实例
func newRouter() *router {
	// 初始化roots映射和handlers映射，并返回router实例
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern解析路由模式字符串，将其拆分为部分片段
// 参数:
//
//	pattern: 路由模式字符串，例如 "/user/:id"
//
// 返回值:
//
//	[]string: 拆分后的路由片段数组
//
// 注意: 如果遇到 "*"，则停止解析，确保只允许一个 "*"
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute为指定的HTTP方法和路由模式添加路由规则及处理函数
// 参数:
//
//	method: HTTP方法，例如 "GET" 或 "POST"
//	pattern: 路由模式字符串，例如 "/user/:id"
//	handler: 处理该路由的函数
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	// 将路由插入到路由树中
	r.roots[method].insert(pattern, parts, 0)
	// 存储处理函数
	r.handlers[key] = handler
}

// getRoute根据HTTP方法和路径查找匹配的路由节点和参数
// 参数:
//
//	method: HTTP方法，例如 "GET" 或 "POST"
//	path: 请求路径，例如 "/user/123"
//
// 返回值:
//
//	*node: 匹配的路由节点，如果没有匹配则返回 nil
//	map[string]string: 路径参数，例如 {"id": "123"}
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	// 在路由树中搜索匹配的节点
	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		// 解析路径参数
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// getRoutes获取指定HTTP方法的所有路由节点
// 参数:
//
//	method: HTTP方法，例如 "GET" 或 "POST"
//
// 返回值:
//
//	[]*node: 路由节点数组，如果没有对应方法的路由则返回 nil
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	// 遍历路由树，收集所有节点
	root.travel(&nodes)
	return nodes
}

// handle处理传入的上下文，找到匹配的路由并调用对应的处理函数
// 参数:
//
//	c: 上下文对象，包含请求和响应信息
func (r *router) handle(c *Context) {
	// 查找匹配的路由和参数
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		// 调用处理函数
		r.handlers[key](c)
	} else {
		// 如果没有匹配的路由，返回404错误
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
