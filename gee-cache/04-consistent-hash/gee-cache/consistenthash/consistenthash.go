package consistenthash

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 定义了一个将字节数据映射为 uint32 的哈希函数类型。
// 该接口允许灵活选择不同的哈希算法。
// 定义了函数类型 Hash，采取依赖注入的方式，允许用于替换成自定义的 Hash 函数，也方便测试时替换，默认为 crc32.ChecksumIEEE 算法。
type Hash func(data []byte) uint32

// Map 表示一个一致性哈希环，包含所有已哈希的键及其对应的节点。
// 它使用一致性哈希算法来分布键到节点上。
// 字段说明：
// - hash: 哈希函数，用于生成哈希值。
// - replicas: 每个节点的虚拟副本数量。
// - keys: 所有哈希值的有序列表。
// - hashMap: 哈希值到实际键的映射。
type Map struct {
	hash     Hash           // 哈希函数
	replicas int            // 虚拟副本数量(虚拟节点倍数)
	keys     []int          // 已排序的哈希值列表(哈希环)
	hashMap  map[int]string // 哈希值到键的映射(虚拟节点与真实节点的映射表)
}

// TODO 构造函数 New() 允许自定义虚拟节点倍数和 Hash 函数。
// New 创建并返回一个新的 Map 实例。
// 参数说明：
// - replicas: 每个节点的虚拟副本数量。
// - fn: 自定义的哈希函数。如果为 nil，则使用默认的 crc32.ChecksumIEEE。
// 返回值：
// - *Map: 一个新的 Map 实例。
func New(replicas int, fn Hash) *Map {
	// 初始化 Map 结构体
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 将一组键添加到一致性哈希环中。
// 参数说明：
// - keys: 需要添加的键列表。
// 功能描述：
// 对于每个键，生成其对应的多个虚拟副本（根据 replicas 数量），并将它们的哈希值加入到 keys 列表中。
// 最后对 keys 列表进行排序以保持有序性。
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 生成虚拟节点的哈希值(“1”+“6” = 16)
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 根据给定的键查找最接近的节点。
// 参数说明：
// - key: 需要查找的键。
// 返回值：
// - string: 最接近该键的节点名称。如果没有可用节点，则返回空字符串。
// 功能描述：
// 计算给定键的哈希值，并通过二分查找找到第一个大于等于该哈希值的键。
// 如果找不到，则返回环中第一个节点（即最小的哈希值对应的节点）。
func (m *Map) Get(key string) string {
	// 如果哈希环为空，则返回空字符串。
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// 使用二分查找定位适当的虚拟副本。
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	fmt.Println(idx % len(m.keys))
	fmt.Println(m.keys[idx%len(m.keys)])
	fmt.Println(m.hashMap[m.keys[idx%len(m.keys)]])
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

func (m *Map) PrintlnTest() {
	fmt.Println("m.keys:", m.keys)
	fmt.Println("m.hashMap:", m.hashMap)
}
