package geecache

import (
	"fmt"
	"gee-web/gee-cache/06-single-flight/gee-cache/singleflight"
	"log"
	"sync"
)

// TODO 负责与外部交互，控制缓存存储和获取的主流程
// Group 是一个缓存命名空间，与分布式加载的数据相关联。
// 它提供从本地缓存或远程对等节点获取数据的能力。
type Group struct {
	name      string     // 缓存组的名称
	getter    Getter     // 数据加载器，用于从外部源获取数据
	mainCache cache      // 主缓存，存储本地缓存数据
	peers     PeerPicker // 对等节点选择器，用于选择远程对等节点
	// use singleflight.Group to make sure that
	// each key is only fetched once
	loader *singleflight.Group
}

// Getter 是一个接口，定义了通过键加载数据的方法。
type Getter interface {
	// Get 根据键从外部源加载数据。
	// 参数:
	//   key - 数据的键
	// 返回值:
	//   []byte - 加载的数据
	//   error - 如果加载失败，则返回错误
	Get(key string) ([]byte, error)
}

// GetterFunc 使用函数实现 Getter 接口。
type GetterFunc func(key string) ([]byte, error)

// Get 实现 Getter 接口方法。
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex              // 读写互斥锁，用于保护 groups 的并发访问
	groups = make(map[string]*Group) // 缓存组的全局注册表
)

// NewGroup 创建一个新的 Group 实例。
// 参数:
//
//	name - 缓存组的名称
//	cacheBytes - 主缓存的最大字节数
//	getter - 数据加载器，用于从外部源获取数据
//
// 返回值:
//
//	*Group - 新创建的缓存组实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup 返回之前通过 NewGroup 创建的指定名称的缓存组。
// 如果不存在该名称的缓存组，则返回 nil。
// 参数:
//
//	name - 缓存组的名称
//
// 返回值:
//
//	*Group - 指定名称的缓存组实例
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get 根据键从缓存中获取值。如果值不在本地缓存中，则尝试从远程对等节点或本地加载器获取。
// 参数:
//
//	key - 数据的键
//
// 返回值:
//
//	ByteView - 缓存中的值
//	error - 如果获取失败，则返回错误
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 尝试从本地缓存获取值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	// 如果本地缓存未命中，则尝试从远程对等节点或本地加载器获取
	return g.load(key)
}

// TODO 实现了 PeerPicker 接口的 HTTPPool 注入到 Group 中。
// RegisterPeers 注册一个 PeerPicker，用于选择远程对等节点。
// 如果已经注册过 PeerPicker，则会触发 panic。
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// TODO 修改 load 方法，使用 PickPeer() 方法选择节点，若非本机节点，则调用 getFromPeer() 从远程获取。若是本机节点或失败，则回退到 getLocally()。
// load 尝试从远程对等节点或本地加载器获取值，并将其填充到本地缓存中。
func (g *Group) load(key string) (value ByteView, err error) {
	// 使用 g.loader.Do 方法，确保每个键只被加载一次
	viewi, err := g.loader.Do(key, func() (interface{}, error) {

		// 如果存在 PeerPicker，则尝试从远程对等节点获取值
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		// 如果远程对等节点获取失败，则从本地加载器获取值
		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

// populateCache 将值填充到本地缓存中。
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// getLocally 从本地加载器获取值，并将其填充到本地缓存中。
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// TODO 实现了 PeerGetter 接口的 httpGetter 从访问远程节点，获取缓存值。
// getFromPeer 从远程对等节点获取值。
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}
