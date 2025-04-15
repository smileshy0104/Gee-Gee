package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 捕获panic，并打印堆栈信息
func trace(message string) string {
	// 获取堆栈信息
	var pcs [32]uintptr
	// 从第3个调用栈开始获取堆栈信息
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	// 遍历堆栈信息
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// Recovery 用于捕获程序运行时panic
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			// 捕获panic
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		// 执行后续处理
		c.Next()
	}
}
