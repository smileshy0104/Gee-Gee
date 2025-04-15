package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
// 缓存 是一个 LRU 缓存。它不是并发安全的。
type Cache struct {
	maxBytes int64                    // 允许使用的最大内存
	nbytes   int64                    // 当前已使用的内存
	ll       *list.List               // 双向链表
	cache    map[string]*list.Element // key是字符串，value是双向链表中对应节点的指针
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value) // 移除缓存中的某条记录时，会调用该回调函数，可以为 nil。
}

// entry 是一个双向链表的节点，它包含了一个键值对，键是字符串，值是一个 Value 类型的值。
type entry struct {
	key   string // 键
	value Value  // 值
}

// Value 是一个接口，它定义了一个 Len 方法，用来返回一个 Value 的长度。
type Value interface {
	Len() int // 返回一个 Value 的长度（用于返回值所占用的内存大小。）
}

// New 实例化 LRU 缓存
// maxBytes 是允许使用的最大内存，onEvicted 是一个回调函数，当缓存中的某条记录被移除时，会调用该回调函数。
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,                       // 允许使用的最大内存
		ll:        list.New(),                     // 双向链表
		cache:     make(map[string]*list.Element), // 缓存
		OnEvicted: onEvicted,                      // 回调函数
	}
}

// Add 向缓存中添加一个键值对。
// 如果键已经存在，则更新该键对应的值，并移动该节点到链表的头部。
// 如果键不存在，则添加一个新的键值对到缓存中，并移动该节点到链表的头部。
// 当缓存使用的内存超过最大限制时，移除最旧的键值对。
func (c *Cache) Add(key string, value Value) {
	// 如果键已经存在，则更新该键对应的值，并移动该节点到链表的头部。
	if ele, ok := c.cache[key]; ok {
		// 更新该键对应的值，并移动该节点到链表的头部。
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 更新已使用的内存大小
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 不存在
		// 添加一个键值对到缓存中
		ele := c.ll.PushFront(&entry{key, value})
		// 将键值对添加到缓存中
		c.cache[key] = ele
		// 更新已使用的内存大小
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 如果设置了最大内存限制，并且当前使用的内存超过了限制，则移除最旧的键值对。
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get 从缓存中获取与指定键关联的值。
// 如果找到了该键，则更新该键对应的值，并移动该节点到链表的头部。
// 如果没有找到该键，则返回零值和false。
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 如果键已经存在，则更新该键对应的值，并移动该节点到链表的头部。
	if ele, ok := c.cache[key]; ok {
		// 更新该键对应的值，并移动该节点到链表的头部。
		c.ll.MoveToFront(ele)
		// 获取节点中的值。
		kv := ele.Value.(*entry)
		// 返回该键对应的值和 true。
		return kv.value, true
	}
	// 如果键不存在，返回零值和false。
	return
}

// TODO 删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）
// RemoveOldest 移除缓存中最旧的条目。
// 当缓存达到其容量限制时，此方法用于释放空间。
// 如果缓存为空，则此方法不执行任何操作。
func (c *Cache) RemoveOldest() {
	// 获取最旧的条目，即双向链表的尾部元素。（取到队首节点，从链表中删除。）
	ele := c.ll.Back()
	// 如果存在最旧的条目，则进行移除操作。
	if ele != nil {
		// 从双向链表中移除最旧的条目。
		c.ll.Remove(ele)
		// 从哈希映射中移除对应的键值对。
		kv := ele.Value.(*entry)
		// 从字典中 c.cache 删除该节点的映射关系。
		delete(c.cache, kv.key)
		// 更新缓存使用的字节数。
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 如果设置了缓存项被移除时的回调函数，则调用该函数。
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
// Len() 用来获取添加了多少条数据。
func (c *Cache) Len() int {
	return c.ll.Len()
}
