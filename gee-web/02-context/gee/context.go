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
	// 响应信息
	StatusCode int // HTTP 响应状态码
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
	}
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
