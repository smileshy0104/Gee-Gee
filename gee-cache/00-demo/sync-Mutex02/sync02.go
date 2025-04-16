package main

import (
	"fmt"
	"sync"
	"time"
)

// TODO Map存在线程安全问题，因为多个 Goroutine 可以同时访问和修改它。（所以需要使用互斥锁sync.mutex）
// m 用于同步互斥访问，以防止并发访问冲突
var m sync.Mutex

// set 用于记录已经打印过的数字，以确保每个数字只打印一次
var set = make(map[int]bool, 0)

// printOnce 确保同一个数字只打印一次，即使在并发环境下也是如此
// num: 待打印的数字
func printOnce(num int) {
	// 上锁以保护共享资源的访问
	m.Lock()
	// 检查数字是否已经存在，如果不存在则打印
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	// 将数字标记为已打印
	set[num] = true
	// 解锁以允许其他协程访问共享资源
	m.Unlock()
}

func main() {
	// 启动10个协程，每个协程都尝试打印数字100
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	// 等待一段时间，以确保所有协程都已完成执行
	time.Sleep(time.Second)
}
