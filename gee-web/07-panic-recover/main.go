package main

import (
	"gee-web/gee-web/07-panic-recover/gee"
	"net/http"
)

func main() {
	// 创建一个新的引擎实例，并使用默认的配置。
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello yyds\n")
	})
	// index out of range for testing Recovery()
	// index超过数组长度，触发panic
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"yyds"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
