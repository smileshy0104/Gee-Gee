package consistenthash

import (
	"fmt"
	"strconv"
	"testing"
)

// TestHashing 测试一致性哈希的正确性。
// 参数：
// t - 测试框架提供的测试对象，用于报告测试结果。
// 功能描述：
// 1. 初始化一个一致性哈希环，虚拟节点数量为3，哈希函数将字符串转换为整数。
// 2. 添加初始键值 "6", "4", "2"，并验证特定键是否映射到正确的节点。
// 3. 添加新键值 "8"，并验证更新后的映射关系是否正确。
func TestHashing(t *testing.T) {
	// 初始化一致性哈希环，虚拟节点数量为3，自定义哈希函数将字符串转换为整数。
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// 添加初始键值 "6", "4", "2"，生成对应的哈希值和虚拟节点。
	hash.Add("6", "4", "2")

	// 定义测试用例，验证特定键是否映射到正确的节点。
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	hash.PrintlnTest()
	// 遍历测试用例，验证每个键是否映射到预期的节点。
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	fmt.Println("====================================================")
	// 添加新键值 "8"，生成对应的哈希值和虚拟节点。
	hash.Add("8")

	// 更新测试用例，验证键 "27" 是否映射到新增的节点 "8"。
	testCases["27"] = "8"

	// 再次遍历测试用例，验证更新后的映射关系是否正确。
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
