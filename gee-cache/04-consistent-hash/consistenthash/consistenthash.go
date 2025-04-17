package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 定义了一个将字节数组映射到 uint32 的函数类型。
// 该类型用于指定哈希算法，允许用户自定义哈希函数。
type Hash func(data []byte) uint32

// Map 是一致性哈希的核心结构，包含所有已哈希的键及其对应的节点。
// keys 字段存储了所有哈希值，并保持排序以支持高效的查找操作。
// hashMap 字段用于存储哈希值与实际键的映射关系。
type Map struct {
	hash     Hash           // 用于计算哈希值的函数
	replicas int            // 每个节点的虚拟副本数量
	keys     []int          // 已排序的哈希值列表
	hashMap  map[int]string // 哈希值到实际键的映射
}

// New 创建一个新的 Map 实例。
// 参数：
//   - replicas: 每个节点的虚拟副本数量，用于提高哈希分布的均匀性。
//   - fn: 自定义的哈希函数，如果为 nil，则默认使用 crc32.ChecksumIEEE。
//
// 返回值：
//   - 返回一个初始化好的 *Map 实例。
func New(replicas int, fn Hash) *Map {
	// 初始化 Map 结构体
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	// 如果未提供自定义哈希函数，则使用默认的 crc32.ChecksumIEEE
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 将一组键添加到一致性哈希环中。
// 参数：
//   - keys: 要添加的键的可变参数列表。
//
// 功能描述：
//
//	对于每个键，根据其虚拟副本数量生成多个哈希值，并将其添加到 keys 和 hashMap 中。
//	最后对 keys 列表进行排序，以便后续高效查找。
func (m *Map) Add(keys ...string) {
	// 遍历每个键，生成虚拟副本并添加到 keys 和 hashMap 中
	for _, key := range keys {
		// 生成虚拟副本
		for i := 0; i < m.replicas; i++ {
			// 计算哈希值（hash：i+key）
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 将哈希值和键添加到 keys 和 hashMap 中
			m.keys = append(m.keys, hash)
			// 将哈希值和键添加到 hashMap 中
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys) // 对哈希值列表进行排序
}

// Get 根据给定的键找到最接近的节点。
// 参数：
//   - key: 要查找的键。
//
// 返回值：
//   - 返回与给定键最接近的节点名称。如果没有可用节点，则返回空字符串。
//
// 功能描述：
//
//	首先计算给定键的哈希值，然后通过二分查找在已排序的 keys 列表中找到第一个大于等于该哈希值的位置。
//	如果未找到匹配项，则循环回到列表的第一个元素。
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	// 计算给定键的哈希值
	hash := int(m.hash([]byte(key)))
	// 使用二分查找定位合适的虚拟副本位置
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 如果未找到匹配项，则循环回到列表的第一个元素
	return m.hashMap[m.keys[idx%len(m.keys)]] // 返回对应的实际键
}
