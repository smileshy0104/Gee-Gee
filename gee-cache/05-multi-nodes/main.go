package main

// 示例API请求和响应
/*
$ curl "http://localhost:9999/api?key=Tom"
630

$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
*/

import (
	"flag"
	"fmt"
	geecache "gee-web/gee-cache/05-multi-nodes/gee-cache"
	"log"
	"net/http"
)

// 示例数据库，存储用户得分
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// createGroup 创建一个geecache组
// 返回值: geecache.Group指针，用于操作缓存
func createGroup() *geecache.Group {
	// 创建一个geecache组，名称为"scores"，最大容量为2KB，GetterFunc为获取数据的函数。
	return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

// TODO startCacheServer() 用来启动缓存服务器：创建 HTTPPool，添加节点信息，注册到 gee 中，启动 HTTP 服务（共3个端口，8001/8002/8003），用户不感知。
// startCacheServer 启动缓存服务器
// 参数:
// - addr: 当前缓存服务器的地址
// - addrs: 所有缓存服务器的地址列表
// - gee: geecache组指针
func startCacheServer(addr string, addrs []string, gee *geecache.Group) {
	peers := geecache.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// TODO startAPIServer() 用来启动一个 API 服务（端口 9999），与用户进行交互，用户感知。
// startAPIServer 启动API服务器
// 参数:
// - apiAddr: API服务器的地址
// - gee: geecache组指针
func startAPIServer(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

// main 是程序的入口点
func main() {
	var port int
	var api bool
	// main() 函数需要命令行传入 port 和 api 2 个参数，用来在指定端口启动 HTTP 服务。
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	// apiAddr 表示 API 服务的地址，默认为 http://localhost:9999
	apiAddr := "http://localhost:9999"
	// addrMap 表示缓存服务器的地址，默认为 http://localhost:8001/8002/8003
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	// addrs 表示所有缓存服务器的地址
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], addrs, gee)
}
