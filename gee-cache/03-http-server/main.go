package main

// 示例中使用curl命令通过HTTP请求访问geecache服务，获取键值对应的值。
// 如果键存在，则返回对应的值，如Tom的值是630。
// 如果键不存在，则返回提示信息，如kkk not exist。
import (
	"fmt"
	geecache "gee-web/gee-cache/03-http-server/gee-cache"
	"log"
	"net/http"
)

// db模拟一个数据库，存储了用户分数信息。
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// main函数是程序的入口点。
// 它初始化了一个geecache组，配置了获取数据的函数，并启动了HTTP服务。
func main() {
	// 创建一个名为"scores"的缓存组，最大容量为2KB。
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		// 定义了当缓存未命中时，如何从数据源获取数据的函数。
		// 如果键在模拟的数据库中存在，则返回对应的值。
		// 如果键不存在，则返回错误信息，表明键不存在。
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	// 定义HTTP服务的地址。
	addr := "localhost:9999"
	// 创建一个HTTPPool，用于管理HTTP缓存服务的节点。
	peers := geecache.NewHTTPPool(addr)
	// 启动HTTP服务，监听指定地址，处理缓存请求。
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
