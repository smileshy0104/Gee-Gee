package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 是一个类型，表示键值对集合，用于存储动态数据。
type H map[string]interface{}

// Context 封装了请求和响应处理的上下文，提供了处理 HTTP 请求和响应的方法。
type Context struct {
	// 原始对象
	Writer http.ResponseWriter // 用于写入响应
	Req    *http.Request       // 保存请求数据
	// 请求信息
	Path   string // 请求路径
	Method string // 请求方法
	Params map[string]string
	// 响应信息
	StatusCode int // HTTP 响应状态码
	// middleware
	handlers []HandlerFunc
	index    int
}

// newContext 创建并返回一个新的 Context 实例。
// 参数:
// - w: http.ResponseWriter，用于写入响应。
// - req: *http.Request，保存请求数据。
// 返回值:
// - *Context: 新创建的 Context 实例。
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next 执行上下文中的下一个处理函数。
// 该方法主要用于在中间件或处理函数链中推进执行顺序。
func (c *Context) Next() {
	// 将索引递增到下一个处理函数的位置。
	c.index++
	// 获取处理函数列表的长度，用于后续的循环条件判断。
	s := len(c.handlers)
	// 遍历剩余的处理函数并执行它们。
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Fail 用于在执行过程中报告失败信息。
// 该方法主要用于中断当前的处理流程，并返回错误信息给客户端。
// 参数:
//
//	code - HTTP状态码，表示错误的类型。
//	err  - 错误消息，提供具体的错误信息。
func (c *Context) Fail(code int, err string) {
	// 将索引设置为处理函数列表的末尾，以中断后续的处理。
	c.index = len(c.handlers)
	// 使用JSON格式返回错误信息和状态码。
	c.JSON(code, H{"message": err})
}

// Param 从 Context 的 Params 映射中获取指定 key 对应的值。
// 如果 key 存在，则返回对应的值；如果 key 不存在，则返回空字符串。
// 参数:
//
//	key - string 类型，表示要获取的参数键。
//
// 返回值:
//
//	string 类型，表示与键关联的值，若键不存在则返回空字符串。
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm 从 POST 表单数据中获取指定 key 的值。
// 参数:
// - key: string，表单字段的键。
// 返回值:
// - string: 对应键的值。
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 从 URL 查询参数中获取指定 key 的值。
// 参数:
// - key: string，查询参数的键。
// 返回值:
// - string: 对应键的值。
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置 HTTP 响应的状态码。
// 参数:
// - code: int，HTTP 状态码。
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头中的指定键值对。
// 参数:
// - key: string，响应头的键。
// - value: string，响应头的值。
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 返回一个带有指定状态码和格式化字符串的纯文本响应。
// 参数:
// - code: int，HTTP 状态码。
// - format: string，格式化字符串模板。
// - values: ...interface{}，格式化字符串的参数列表。
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 返回一个带有指定状态码和 JSON 数据的响应。
// 参数:
// - code: int，HTTP 状态码。
// - obj: interface{}，要序列化的对象。
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 返回一个带有指定状态码和字节数组数据的响应。
// 参数:
// - code: int，HTTP 状态码。
// - data: []byte，要发送的字节数组数据。
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 返回一个带有指定状态码和 HTML 内容的响应。
// 参数:
// - code: int，HTTP 状态码。
// - html: string，HTML 内容字符串。
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
