package geecache

// TODO 缓存值的抽象与封装
// ByteView 表示一个不可变的字节视图。
type ByteView struct {
	b []byte // 存储实际的字节数据
}

// Len 返回 ByteView 中字节切片的长度。
//
// 返回值:
//
//	int - 字节切片的长度。
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回 ByteView 中数据的一个副本，以字节切片的形式。
//
// 返回值:
//
//	[]byte - 数据的副本，确保原始数据不会被修改。
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 将 ByteView 中的数据以字符串形式返回。
// 如果需要，会创建数据的副本。
//
// 返回值:
//
//	string - 数据的字符串表示。
func (v ByteView) String() string {
	return string(v.b)
}

// cloneBytes 创建并返回字节切片 b 的副本。
//
// 参数:
//
//	b []byte - 需要复制的字节切片。
//
// 返回值:
//
//	[]byte - 输入字节切片的副本。
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b)) // 创建与输入切片等长的新切片
	copy(c, b)                // 将输入切片的内容复制到新切片中
	return c                  // 返回副本
}
