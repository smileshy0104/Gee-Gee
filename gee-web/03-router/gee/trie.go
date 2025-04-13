package gee

import (
	"fmt"
	"strings"
)

// node 是一个树结构的节点，用于存储路由的路径信息。
// pattern: 完整的路由路径。
// part: 路由路径的一部分。
// children: 子节点。
// isWild: 是否是动态参数（以 ':' 或 '*' 开头）。
type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

// String 实现了 fmt.Stringer 接口，用于打印节点信息。
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// insert 递归地在节点中插入一个路由路径。
// pattern: 完整的路由路径。
// parts: 路由路径的各个部分。
// height: 当前处理的路径部分的索引。
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search 递归地在节点中搜索一个路由路径。
// parts: 路由路径的各个部分。
// height: 当前处理的路径部分的索引。
// 返回值: 如果找到完整的路由路径，则返回对应的节点；否则返回 nil。
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// travel 遍历节点树，收集所有有完整路由路径的节点。
// list: 用于存储找到的所有节点的切片指针。
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// matchChild 查找与指定路径部分匹配的子节点。
// part: 路径部分。
// 返回值: 如果找到匹配的子节点，则返回它；否则返回 nil。
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 查找与指定路径部分匹配的所有子节点。
// part: 路径部分。
// 返回值: 包含所有匹配子节点的切片。
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
