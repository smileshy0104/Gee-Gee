package geecache

import (
	"fmt"
	"gee-web/gee-cache/05-multi-nodes/gee-cache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// TODO 提供被其他节点访问的能力(基于http)
const (
	defaultBasePath = "/_geecache/" // 默认的基础路径
	defaultReplicas = 50            // 一致性哈希的副本数量
)

// HTTPPool 实现了 PeerPicker 接口，用于管理一组 HTTP 节点（peers）。
type HTTPPool struct {
	self     string              // 当前节点的基本 URL，例如 "https://example.net:8000"
	basePath string              // 基础路径
	mu       sync.Mutex          // 保护 peers 和 httpGetters 的互斥锁
	peers    *consistenthash.Map // 一致性哈希映射的节点（是一致性哈希算法的 Map，用来根据具体的 key 选择节点。）
	// 每一个远程节点对应一个 httpGetter，因为 httpGetter 与远程节点的地址 baseURL 有关。
	httpGetters map[string]*httpGetter // HTTP Getter，按 URL 键入，例如 "http://10.0.0.2:8008"
}

// NewHTTPPool 初始化一个 HTTP 节点池。
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log 记录以服务器名称为前缀的信息
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP 处理所有 HTTP 请求
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 检查请求路径是否以 basePath 开头
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// 解析请求路径，期望格式为 /<basepath>/<groupname>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0] // 组名
	key := parts[1]       // 键

	group := GetGroup(groupName) // 获取对应的缓存组
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key) // 从组中获取数据
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream") // 设置响应头
	w.Write(view.ByteSlice())                                  // 写入响应体
}

// Set() 方法实例化了一致性哈希算法，并且添加了传入的节点。
// Set 更新节点池的节点列表。
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()                                              // 加锁以保护共享资源
	defer p.mu.Unlock()                                      // 解锁
	p.peers = consistenthash.New(defaultReplicas, nil)       // 创建一致性哈希映射
	p.peers.Add(peers...)                                    // 添加节点
	p.httpGetters = make(map[string]*httpGetter, len(peers)) // 初始化 HTTP Getter 映射
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath} // 并为每一个节点创建了一个 HTTP 客户端 httpGetter。
	}
}

// PickPeer 根据键选择一个节点
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()                                                 // 加锁
	defer p.mu.Unlock()                                         // 解锁
	if peer := p.peers.Get(key); peer != "" && peer != p.self { // 获取节点
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true // 返回对应的 HTTP Getter
	}
	return nil, false // 未找到合适的节点
}

var _ PeerPicker = (*HTTPPool)(nil) // 确保 HTTPPool 实现了 PeerPicker 接口

// httpGetter 结构体用于通过 HTTP 获取数据
type httpGetter struct {
	// baseURL 表示将要访问的远程节点的地址，例如 http://example.com/_geecache/。
	baseURL string // 基本 URL
}

// PickerPeer() 包装了一致性哈希算法的 Get() 方法，根据具体的 key，选择节点，返回节点对应的 HTTP 客户端。
// Get 从指定组和键获取数据
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",              // 构建请求 URL
		h.baseURL,              // baseURL 表示将要访问的远程节点的地址，例如 http://example.com/_geecache/。
		url.QueryEscape(group), // 对组名进行 URL 编码
		url.QueryEscape(key),   // 对键进行 URL 编码
	)
	// 使用 http.Get() 方式获取返回值，并转换为 []bytes 类型。
	res, err := http.Get(u) // 发送 HTTP GET 请求
	if err != nil {
		return nil, err // 返回错误
	}
	defer res.Body.Close() // 确保在函数结束时关闭响应体

	if res.StatusCode != http.StatusOK { // 检查响应状态
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body) // 读取响应体
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil // 返回获取到的数据
}

var _ PeerGetter = (*httpGetter)(nil) // 确保 httpGetter 实现了 PeerGetter 接口
