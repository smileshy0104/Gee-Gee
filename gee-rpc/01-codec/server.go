package geerpc

import (
	"encoding/json"
	"fmt"
	"gee-web/gee-rpc/01-codec/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

// MagicNumber 是用于标识 geerpc 请求的魔数。
const MagicNumber = 0x3bef5c

// Option 定义了 RPC 请求的选项。
type Option struct {
	MagicNumber int        // 魔数，用于标记这是一个 geerpc 请求
	CodecType   codec.Type // 客户端可以选择不同的编解码器来编码请求体
}

// DefaultOption 是默认的 Option 实例。
var DefaultOption = &Option{
	MagicNumber: MagicNumber,   // 魔数，用于标记这是一个 geerpc 请求
	CodecType:   codec.GobType, // 客户端可以选择不同的编解码器来编码请求体
}

// Server 表示一个 RPC 服务器。
type Server struct{}

// NewServer 返回一个新的 Server 实例。
func NewServer() *Server {
	return &Server{}
}

// DefaultServer 是默认的 Server 实例。
var DefaultServer = NewServer()

// ServeConn 在单个连接上运行服务器。
// ServeConn 会阻塞，直到客户端断开连接。
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	// 确保连接被关闭
	defer func() { _ = conn.Close() }()
	var opt Option
	// 从连接中读取选项
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	// 检查选项的魔数
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	// 根据选项的编解码器类型创建编解码器
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	// 处理连接
	server.serveCodec(f(conn))
}

// invalidRequest 是当发生错误时响应参数的占位符。
var invalidRequest = struct{}{}

// serveCodec 处理通过指定编解码器发送的请求。
func (server *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex) // 确保发送完整的响应
	wg := new(sync.WaitGroup)  // 等待所有请求处理完成
	for {
		// 读取请求
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break // 无法恢复，关闭连接
			}
			req.h.Error = err.Error()
			// 发送错误响应
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		// 处理请求
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

// request 存储调用的所有信息。
type request struct {
	h      *codec.Header // 请求头
	argv   reflect.Value // 请求参数
	replyv reflect.Value // 响应值
}

// readRequestHeader 从编解码器中读取请求头。
func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	// 从连接中读取请求头。
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

// readRequest 从编解码器中读取完整的请求。
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	// 读取请求头
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	// TODO: 目前我们不知道请求参数的类型
	// 第一天，假设它是字符串
	req.argv = reflect.New(reflect.TypeOf(""))
	// 从连接中读取请求参数
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

// sendResponse 使用指定的编解码器发送响应。
func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	// 使用指定的编解码器发送响应
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

// handleRequest 处理单个请求。
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// TODO: 应该调用注册的 RPC 方法以获取正确的响应值
	// 第一天，只是打印参数并发送一条问候消息
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Seq))
	// 发送响应
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

// Accept 接受监听器上的连接并为每个传入连接提供服务。
func (server *Server) Accept(lis net.Listener) {
	for {
		// 接受客户端的连接请求
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		// 为每个连接提供服务
		go server.ServeConn(conn)
	}
}

// Accept 接受监听器上的连接并为每个传入连接提供服务。
func Accept(lis net.Listener) {
	// 将服务器绑定地址发送到addr通道
	DefaultServer.Accept(lis)
}
