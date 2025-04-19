package singleflight

import "sync"

// call 表示正在进行中或已完成的 Do 调用。
// 它包含函数调用的结果、错误以及用于同步的 WaitGroup。
type call struct {
	wg  sync.WaitGroup // 用于等待调用完成
	val interface{}    // 函数调用的结果（使用 sync.WaitGroup 锁避免重入。）
	err error          // 函数调用的错误
}

// Group 是 singleflight 的主数据结构，管理不同 key 的请求(call)。
// Group 表示一类工作，并形成一个命名空间，
// 在该命名空间中可以执行任务并实现重复请求的抑制。
type Group struct {
	mu sync.Mutex       // 保护对 m 的并发访问
	m  map[string]*call // 懒初始化的映射，存储正在进行中的调用
}

// TODO Do 的作用就是，针对相同的 key，无论 Do 被调用多少次，函数 fn 都只会被调用一次，等待 fn 调用结束了，返回返回值或错误。
// Do 执行并返回给定函数 fn 的结果，确保对于相同的 key，
// 同一时间只有一个正在执行的调用。如果有重复请求到来，
// 重复的调用者会等待原始调用完成，并接收相同的结果。
//
// 参数:
//   - key: string 类型，表示调用的唯一标识符，用于检测重复请求。
//   - fn: func() (interface{}, error) 类型，表示实际需要执行的函数。
//
// 返回值:
//   - interface{}: 函数 fn 的执行结果。
//   - error: 函数 fn 执行过程中可能产生的错误。
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	// 如果映射 m 尚未初始化，则进行懒初始化。
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 如果存在相同 key 的调用，说明是重复请求。
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// 等待之前的调用完成。
		c.wg.Wait()
		return c.val, c.err
	}
	// 创建新的 call 实例并加入映射。
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 执行用户提供的函数 fn，并记录结果和错误。
	c.val, c.err = fn()
	// 标记调用完成，唤醒等待的调用者。
	c.wg.Done()

	// 清理已完成的调用。
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	// 返回函数调用的结果和错误。
	return c.val, c.err
}
