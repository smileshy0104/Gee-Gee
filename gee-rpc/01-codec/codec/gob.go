package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

// GobCodec 实现了 Codec 接口，使用 gob 包进行序列化和反序列化。
// 它负责通过网络连接对数据进行编码和解码。
type GobCodec struct {
	conn io.ReadWriteCloser // 网络连接，用于读写数据。
	buf  *bufio.Writer      // 缓冲写入器，用于优化写操作。
	dec  *gob.Decoder       // gob 解码器，用于从连接中解码数据。
	enc  *gob.Encoder       // gob 编码器，用于将数据编码后写入连接。
}

var _ Codec = (*GobCodec)(nil)

// NewGobCodec 创建一个新的 GobCodec 实例。
// 参数:
//   - conn: io.ReadWriteCloser 类型的网络连接，用于数据传输。
//
// 返回值:
//   - Codec 接口的实现，用于处理 gob 编码和解码。
func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

// ReadHeader 从连接中读取并解码 Header 数据。
// 参数:
//   - h: 指向 Header 的指针，用于存储解码后的头部信息。
//
// 返回值:
//   - error: 如果解码失败，则返回错误；否则返回 nil。
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

// ReadBody 从连接中读取并解码消息体数据。
// 参数:
//   - body: 接收解码后数据的接口，通常是一个结构体指针。
//
// 返回值:
//   - error: 如果解码失败，则返回错误；否则返回 nil。
func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

// Write 将 Header 和消息体数据编码后写入连接。
// 参数:
//   - h: 指向 Header 的指针，包含需要编码的头部信息。
//   - body: 需要编码的消息体数据，通常是一个结构体。
//
// 返回值:
//   - error: 如果编码或写入失败，则返回错误；否则返回 nil。
//
// 注意: 在函数结束时会刷新缓冲区，并在发生错误时关闭连接。
func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush() // 刷新缓冲区以确保所有数据被写入连接。
		if err != nil {
			_ = c.Close() // 如果发生错误，则关闭连接以释放资源。
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc: gob error encoding header:", err)
		return
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc: gob error encoding body:", err)
		return
	}
	return
}

// Close 关闭底层的网络连接。
// 返回值:
//   - error: 如果关闭连接失败，则返回错误；否则返回 nil。
func (c *GobCodec) Close() error {
	return c.conn.Close()
}
