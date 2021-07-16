package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由 example /p/:lang
	part     string  // 路由中的一部分
	children []*node // 字节点 [doc, tutorial, intro]
	isWild   bool    //是否精确匹配, part 含有 : 或者 * 为 true
}


// matchChild 匹配第一个成功的节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 所有匹配成功的节点, 用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part ==part || child.isWild {
			nodes = append(nodes,child)
		}
	}
	return nodes
}

// insert 树上插入一个新节点
func (n *node)insert(pattern string, parts []string, height int)  {
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

// search 树上搜索匹配的节点
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

