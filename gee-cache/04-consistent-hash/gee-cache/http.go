package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// defaultBasePath 是 HTTP 请求的默认基础路径。
const defaultBasePath = "/_geecache/"

// HTTPPool 实现了 PeerPicker 接口，用于管理一组通过 HTTP 通信的对等节点。
type HTTPPool struct {
	// self 是当前节点的基础 URL，例如 "https://example.net:8000"。
	self string
	// basePath 是 HTTP 请求的基础路径。
	basePath string
}

// NewHTTPPool 初始化一个新的 HTTPPool 实例。
// 参数：
//   - self: 当前节点的基础 URL。
//
// 返回值：
//   - *HTTPPool: 返回初始化后的 HTTPPool 指针。
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log 记录带有服务器名称的日志信息。
// 参数：
//   - format: 日志格式化字符串。
//   - v: 可变参数，用于格式化日志内容。
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP 处理所有的 HTTP 请求。
// 参数：
//   - w: http.ResponseWriter，用于写入 HTTP 响应。
//   - r: *http.Request，表示接收到的 HTTP 请求。
//
// 功能描述：
//   - 验证请求路径是否以 basePath 开头，如果不是则抛出异常。
//   - 解析请求路径为 <groupname> 和 <key> 两部分。
//   - 根据 groupName 获取对应的缓存组，如果不存在则返回 404 错误。
//   - 使用 group.Get(key) 获取缓存数据，如果发生错误则返回 500 错误。
//   - 将缓存数据以二进制流的形式返回给客户端。
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 如果请求路径不以 basePath 开头，则抛出异常。
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	// 记录请求方法和路径的日志信息。
	p.Log("%s %s", r.Method, r.URL.Path)

	// 解析请求路径为 <groupname> 和 <key> 两部分。
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	// 根据 groupName 获取缓存组，如果不存在则返回 404 错误。
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	// 获取缓存数据，如果发生错误则返回 500 错误。
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头为二进制流类型，并将缓存数据写入响应体。
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
