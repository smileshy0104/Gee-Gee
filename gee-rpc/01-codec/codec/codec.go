package codec

import "io"

// Header 结构体用于定义消息的头部信息。
type Header struct {
	ServiceMethod string // ServiceMethod 是服务名和方法名，通常与 Go 语言中的结构体和方法相映射。
	Seq           uint64 // Seq 是请求的序号，也可以认为是某个请求的 ID，用来区分不同的请求。
	Error         string // Error 是错误信息，客户端置为空，服务端如果如果发生错误，将错误信息置于 Error 中。
}

// Codec 接口定义了网络编解码器需要实现的基本方法。
type Codec interface {
	io.Closer                         // Codec 继承了 io.Closer 接口，确保编解码器可以被正确关闭。
	ReadHeader(*Header) error         // ReadHeader 方法用于读取并解析消息头部。
	ReadBody(interface{}) error       // ReadBody 方法用于读取并解析消息体。
	Write(*Header, interface{}) error // Write 方法用于写入消息，包括头部和体。
}

// NewCodecFunc 类型用于定义创建编解码器实例的函数。
type NewCodecFunc func(io.ReadWriteCloser) Codec

// Type 类型用于定义编解码器的类型。
type Type string

// 定义编解码器类型的常量。
const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

// NewCodecFuncMap 是一个映射，用于根据编解码器类型快速查找对应的创建函数。
var NewCodecFuncMap map[Type]NewCodecFunc

// init 函数用于初始化 NewCodecFuncMap，并注册已有的编解码器创建函数。
func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
