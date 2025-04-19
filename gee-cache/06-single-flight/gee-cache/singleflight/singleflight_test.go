package singleflight

import (
	"testing"
)

// TestDo 测试 Group 的 Do 方法是否能正确执行指定的键值操作并返回预期结果。
// 该测试验证了 Do 方法在处理一个简单的、无错误的操作时的表现。
func TestDo(t *testing.T) {
	// 初始化一个 Group 实例变量 g，用于调用 Do 方法。
	var g Group

	// 调用 Do 方法，指定键为"key"，操作为返回字符串"bar"。
	// 该调用应返回结果"bar"和nil错误。
	v, err := g.Do("key", func() (interface{}, error) {
		return "bar", nil
	})

	// 检查 Do 方法的返回值和错误是否与预期相符。
	// 如果结果不为"bar"或存在错误，则测试失败。
	if v != "bar" || err != nil {
		t.Errorf("Do v = %v, error = %v", v, err)
	}
}
