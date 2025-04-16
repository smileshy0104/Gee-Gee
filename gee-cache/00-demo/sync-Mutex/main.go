// main 包定义了程序的入口点。
package main

// 导入必要的包
import (
	"fmt"
	"time"
)

// TODO Map存在线程安全问题，因为多个 Goroutine 可以同时访问和修改它。（所以需要使用互斥锁sync.mutex）
// set 是一个全局变量，用于记录已经打印过的数字，确保每个数字只打印一次。
var set = make(map[int]bool, 0)

// printOnce 函数确保传入的数字仅被打印一次。
// 参数:
//
//	num - 需要打印的整数。
func printOnce(num int) {
	// 检查该数字是否已经存在于 set 中。如果不存在，则打印该数字。
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	// 将该数字标记为已打印。
	set[num] = true
}

func main() {
	// 启动 10 个 Goroutine，每个 Goroutine 调用 printOnce 函数尝试打印数字 100。
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	// 主线程休眠 1 秒钟，确保所有 Goroutine 有足够的时间执行完成。
	time.Sleep(time.Second)
}
