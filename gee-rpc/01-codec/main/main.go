package main

import (
	"encoding/json"
	"fmt"
	geerpc "gee-web/gee-rpc/01-codec"
	"gee-web/gee-rpc/01-codec/codec"
	"log"
	"net"
	"time"
)

// TODO 在这里我们已经实现了一个消息的编解码器 GobCodec，并且客户端与服务端实现了简单的协议交换(protocol exchange)，即允许客户端使用不同的编码方式。同时实现了服务端的雏形，建立连接，读取、处理并回复客户端的请求。
// startServer 启动一个RPC服务器，监听一个随机端口，并将绑定的地址通过addr通道发送给调用方。
// 参数：
//
//	addr - 用于传递服务器绑定地址的通道。
func startServer(addr chan string) {
	// 监听一个随机可用的TCP端口
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

// main 函数是程序的入口点，启动一个RPC服务器并模拟一个简单的客户端与之通信。
func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	// 模拟一个简单的geerpc客户端，连接到刚启动的服务器
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// 向服务器发送默认选项
	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// 循环发送请求并接收响应
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}
