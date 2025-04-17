package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

// db模拟一个数据库，存储了一些用户的数据
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// TestGetter 测试 GetterFunc 类型的回调函数是否正确工作
func TestGetter(t *testing.T) {
	// 定义一个 GetterFunc 类型的变量 f（回调函数）
	// 借助 GetterFunc 的类型转换，将一个匿名回调函数转换成了接口 f Getter。
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		fmt.Println("getterFunc，我成功打印！")
		return []byte(key), nil
	})

	// 期望得到的结果
	expect := []byte("key")
	// 使用 f.Get 方法获取结果，并与期望值进行比较
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("callback failed")
	}
}

// TestGet 测试 gee.Get 方法是否能正确从缓存或数据库中获取数据
func TestGet(t *testing.T) {
	// 记录每个键从数据库加载的次数
	loadCounts := make(map[string]int, len(db))
	// 创建一个名为 "scores" 的缓存组
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			// 检查键是否存在于数据库中
			if v, ok := db[key]; ok {
				// 记录键的加载次数
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	// 遍历数据库中的每个键值对，测试缓存是否能正确获取值
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		// 再次获取值，以测试缓存是否正常工作，以及负载计数是否正确
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	// 测试当请求一个不存在的键时，缓存是否返回错误
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

// TestGetGroup 测试 GetGroup 函数是否能正确返回指定名称的缓存组
func TestGetGroup(t *testing.T) {
	groupName := "scores"
	// 创建一个缓存组，尽管回调函数为空，但足以测试 GetGroup 函数
	NewGroup(groupName, 2<<10, GetterFunc(
		func(key string) (bytes []byte, err error) {
			fmt.Println("getterFunc，我成功打印！")
			return
		}))
	// 测试 GetGroup 是否能返回正确名称的缓存组
	if group := GetGroup(groupName); group == nil || group.name != groupName {
		t.Fatalf("group %s not exist", groupName)
	}

	// 测试当组不存在时，GetGroup 是否返回 nil
	if group := GetGroup(groupName + "111"); group != nil {
		t.Fatalf("expect nil, but %s got", group.name)
	}
}
