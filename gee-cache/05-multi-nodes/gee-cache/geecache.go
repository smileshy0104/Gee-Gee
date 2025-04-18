package geecache

import (
	"fmt"
	"log"
	"sync"
)

// TODO 负责与外部交互，控制缓存存储和获取的主流程
// Group 是一个缓存命名空间及其关联的数据加载逻辑。
// 它管理一个命名空间和该命名空间的数据加载逻辑。
type Group struct {
	name      string // 缓存组的名称
	getter    Getter // 数据加载器接口
	mainCache cache  // 主缓存实例
}

// Getter 是一个接口，用于根据键加载数据。
// Get 方法接收一个键并返回对应的字节数据和错误信息。
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 是一个函数类型，实现了 Getter 接口。
type GetterFunc func(key string) ([]byte, error)

// Get 实现了 Getter 接口的方法。
// 它调用底层的函数来获取数据。
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 全局变量 mu 用于保护 groups 的并发访问。
var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) // 存储所有创建的 Group 实例
)

// NewGroup 创建一个新的 Group 实例。
// 参数：
//   - name: 缓存组的名称
//   - cacheBytes: 缓存的最大容量（以字节为单位）
//   - getter: 数据加载器
//
// 返回值：
//   - *Group: 新创建的 Group 实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter") // 如果 getter 为 nil，则抛出异常
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g // 将新创建的 Group 实例存储到全局 map 中
	return g
}

// GetGroup 根据名称获取之前创建的 Group 实例。
// 参数：
//   - name: 缓存组的名称
//
// 返回值：
//   - *Group: 找到的 Group 实例，如果不存在则返回 nil
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get 从缓存中获取指定键的值。
// 参数：
//   - key: 缓存键
//
// 返回值：
//   - ByteView: 缓存值的视图
//   - error: 错误信息，如果键为空或加载失败则返回错误
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required") // 如果键为空，返回错误
	}

	// 检查主缓存中是否存在该键
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit") // 如果命中缓存，记录日志
		return v, nil
	}

	// 如果未命中缓存，则尝试加载数据
	return g.load(key)
}

// load 加载指定键的值。
// 参数：
//   - key: 缓存键
//
// 返回值：
//   - ByteView: 加载的缓存值视图
//   - error: 错误信息，如果加载失败则返回错误
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key) // 调用本地加载方法
}

// getLocally 从本地加载指定键的值，并将其添加到缓存中。
// 参数：
//   - key: 缓存键
//
// 返回值：
//   - ByteView: 加载的缓存值视图
//   - error: 错误信息，如果加载失败则返回错误
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key) // 调用 Getter 加载数据
	if err != nil {
		return ByteView{}, err // 如果加载失败，返回错误
	}

	value := ByteView{b: cloneBytes(bytes)} // 克隆字节数组以避免修改原始数据
	g.populateCache(key, value)             // 将加载的数据添加到缓存中
	return value, nil
}

// populateCache 将指定键值对添加到主缓存中。
// 参数：
//   - key: 缓存键
//   - value: 缓存值视图
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value) // 将键值对添加到主缓存中
}
