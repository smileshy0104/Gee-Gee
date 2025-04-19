package geecache

// TODO　PeerPicker 的 PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGetter。
// PeerPicker 是一个接口，用于实现定位拥有特定 key 的 peer 的功能。
//
// 方法说明：
// PickPeer 根据给定的 key 选择对应的 peer。
// 参数：
//   - key: 用于定位的键值。
//
// 返回值：
//   - peer: 实现了 PeerGetter 接口的 peer 对象，用于获取缓存数据。
//   - ok: 表示是否成功找到对应的 peer。
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerPickerFunc func(key string) (peer PeerGetter, ok bool)

// TODO  PeerGetter 的 Get() 方法用于从对应 group 查找缓存值。
// PeerGetter 是一个接口，定义了 peer 必须实现的数据获取方法。
//
// 方法说明：
// Get 根据 group 和 key 获取对应的缓存数据。
// 参数：
//   - group: 缓存组的名称。
//   - key: 缓存项的键值。
//
// 返回值：
//   - []byte: 缓存数据的内容，如果未找到则可能返回空切片。
//   - error: 如果发生错误，则返回具体的错误信息；否则返回 nil。
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
