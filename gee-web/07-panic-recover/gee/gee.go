package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc 定义了gee使用的请求处理函数。
// 它接受一个指向Context的指针作为参数，允许访问请求上下文。
type HandlerFunc func(*Context)

// Engine 实现了ServeHTTP接口。
// 它管理HTTP请求的路由和中间件。
type (
	RouterGroup struct {
		prefix      string        // 该组的URL路径前缀。
		middlewares []HandlerFunc // 该组的中间件函数。
		parent      *RouterGroup  // 父组，支持嵌套。
		engine      *Engine       // 引用Engine实例，所有组共享。
	}

	Engine struct {
		*RouterGroup                     // 嵌入的RouterGroup，用于Engine。
		router        *router            // 用于处理请求路由的路由器。
		groups        []*RouterGroup     // 存储所有RouterGroup。
		htmlTemplates *template.Template // for html render
		funcMap       template.FuncMap   // for html render
	}
)

// New 是gee.Engine的构造函数。
// 它初始化一个新的Engine实例，带有新的路由器和默认的RouterGroup。
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Default 创建一个新的Engine实例，并添加默认的中间件。
// 返回:
//   - *Engine: 创建的Engine实例。
func Default() *Engine {
	// 创建一个新的Engine实例，并添加默认的中间件。
	engine := New()
	// 添加默认的中间件。
	engine.Use(Logger(), Recovery())
	return engine
}

// Group 用于创建一个新的RouterGroup。
// 记住所有组共享相同的Engine实例。
// 参数:
//   - prefix: 该组的URL路径前缀。
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use 用于添加中间件。
// 参数:
//   - middlewares: 中间件函数。
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// addRoute 用于添加路由。
// 参数:
//   - method: HTTP方法（如GET、POST）。
//   - comp: 路径组件。
//   - handler: 处理该路由的HandlerFunc。
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 用于添加GET请求。
// 参数:
//   - pattern: 请求路径模式。
//   - handler: 处理该请求的HandlerFunc。
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 用于添加POST请求。
// 参数:
//   - pattern: 请求路径模式。
//   - handler: 处理该请求的HandlerFunc。
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// createStaticHandler 创建一个处理静态文件的HandlerFunc。
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static 用于处理静态文件。
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// SetFuncMap 用于设置模板函数。
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 用于加载HTML模板。
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// Run 用于启动HTTP服务器。
// 参数:
//   - addr: 监听地址。
//
// 返回:
//   - err: 启动过程中可能发生的错误。
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 实现了ServeHTTP接口。
// 它处理传入的HTTP请求。
// 参数:
//   - w: HTTP响应写入器。
//   - req: HTTP请求。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
