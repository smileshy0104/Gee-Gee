package geecache

import (
	"gee-web/gee-cache/02-single-node/gee-cache/lru"
	"sync"
)

// TODO 并发控制
// cache 是一个固定大小的缓存，使用LRU(最近最少使用)策略来移除缓存项。
type cache struct {
	mu         sync.Mutex // 互斥锁，用于确保并发安全
	lru        *lru.Cache // LRU缓存实例
	cacheBytes int64      // 缓存的总字节大小限制
}

// add 方法向缓存中添加一项。
// 参数:
//
//	key: 用于标识缓存项的键。
//	value: 要添加到缓存的值，类型为 ByteView。
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 如果缓存为空，则创建一个新的LRU缓存实例
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	// 将键值对添加到LRU缓存中
	c.lru.Add(key, value)
}

// get 方法从缓存中获取一项。
// 参数:
//
//	key: 要检索的缓存项的键。
//
// 返回值:
//
//	value: 如果找到缓存项，则返回该项的值，类型为 ByteView。
//	ok: 布尔值，表示是否成功找到缓存项。
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	// 从LRU缓存中获取键对应的值
	if v, ok := c.lru.Get(key); ok {
		// 将值转换为ByteView类型并返回
		return v.(ByteView), ok
	}

	return
}
