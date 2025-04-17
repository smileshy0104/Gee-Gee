package consistenthash

import (
	"strconv"
	"testing"
)

// TestHashing 测试一致性哈希的正确性。
// 参数：
// - t *testing.T: 测试框架提供的测试上下文对象，用于报告测试结果。
// 功能描述：
// 该函数通过创建一个一致性哈希环并添加节点，验证不同键是否能正确映射到预期的节点。
func TestHashing(t *testing.T) {
	// 创建一个新的一致性哈希环，设置副本数为3，并使用自定义哈希函数。
	// 自定义哈希函数将键转换为整数并返回其哈希值。
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// 向哈希环中添加节点 "6", "4", "2"。
	// 这些节点将在哈希环中生成对应的副本哈希值。
	hash.Add("6", "4", "2")

	// 定义测试用例，包含键及其期望映射的节点。
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	// 遍历测试用例，验证每个键是否正确映射到预期的节点。
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// 向哈希环中添加新节点 "8"。
	hash.Add("8")

	// 更新测试用例，键 "27" 应该映射到新添加的节点 "8"。
	testCases["27"] = "8"

	// 再次遍历测试用例，验证每个键是否正确映射到更新后的预期节点。
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
