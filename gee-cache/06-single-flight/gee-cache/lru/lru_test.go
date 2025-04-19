package lru

import (
	"reflect"
	"testing"
)

// String 是一个自定义类型，表示字符串，并实现了 Len() 方法以返回字符串长度。
type String string

// Len 返回字符串的长度。
func (d String) Len() int {
	return len(d)
}

// TestGet 测试 LRU 缓存的 Get 方法是否能够正确获取缓存中的值。
// 参数：
//
//	t - *testing.T 类型，用于测试框架的断言和错误报告。
//
// 逻辑：
//  1. 创建一个无容量限制的 LRU 缓存实例。
//  2. 添加键值对 "key1" -> "1234" 到缓存中。
//  3. 验证通过 Get 方法获取的值是否与预期一致。
//  4. 验证不存在的键 "key2" 是否返回未命中。
func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

// TestRemoveoldest 测试 LRU 缓存的 RemoveOldest 方法是否能够正确移除最旧的缓存项。
// 参数：
//
//	t - *testing.T 类型，用于测试框架的断言和错误报告。
//
// 逻辑：
//  1. 计算缓存容量，确保添加多个键值对后触发移除最旧项。
//  2. 添加三个键值对到缓存中。
//  3. 验证最旧的键 "key1" 是否已被移除，且缓存大小是否符合预期。
func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

// TestOnEvicted 测试 LRU 缓存的 OnEvicted 回调函数是否在移除项时被正确调用。
// 参数：
//
//	t - *testing.T 类型，用于测试框架的断言和错误报告。
//
// 逻辑：
//  1. 定义一个回调函数，记录被移除的键。
//  2. 创建一个具有固定容量的 LRU 缓存实例，并设置回调函数。
//  3. 添加多个键值对到缓存中，触发移除最旧项。
//  4. 验证回调函数记录的被移除键是否与预期一致。
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

// TestAdd 测试 LRU 缓存的 Add 方法是否能够正确更新键值对并维护字节计数。
// 参数：
//
//	t - *testing.T 类型，用于测试框架的断言和错误报告。
//
// 逻辑：
//  1. 创建一个无容量限制的 LRU 缓存实例。
//  2. 添加两个相同的键值对，覆盖原有值。
//  3. 验证缓存的字节计数是否正确反映最新值的大小。
func TestAdd(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key", String("1"))
	lru.Add("key", String("111"))

	if lru.nbytes != int64(len("key")+len("111")) {
		t.Fatal("expected 6 but got", lru.nbytes)
	}
}
