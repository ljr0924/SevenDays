package my_http

import "strings"

type node struct {
    Pattern  string  // 待匹配路由
    Part     string  // 路由中的一部分 例如 :lang
    Children []*node // 子节点，例如：[doc, tutorial, intro]
    IsWild   bool    // 是否精确匹配 part含有 : 或者 * 时为true
}


// 第一个匹配成功的节点， 用于插入
func (n *node) matchChild(part string) *node {
    for _, child := range n.Children {
        if child.Part == part || child.IsWild {
            return child
        }
    }
    return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
    nodes := make([]*node, 0)
    for _, child := range n.Children {
        if child.Part == part || child.IsWild {
            nodes = append(nodes, child)
        }
    }
    return nodes
}

func (n *node) insert(pattern string, parts[]string, height int) {
    if len(parts) == height {
        n.Pattern = pattern
        return
    }

    part := parts[height]
    child := n.matchChild(part)
    if child == nil {
        child = &node{
            Part:   part,
            IsWild: part[0] == ':' || part == "*",
        }
        n.Children = append(n.Children, child)
    }
    child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
    if len(parts) == height || strings.HasPrefix(n.Part, "*") {
        if n.Pattern == "" {
            return nil
        }
        return n
    }

    part := parts[height]
    children := n.matchChildren(part)

    for _, c := range children {
        result := c.search(parts, height+1)
        if result != nil {
            return result
        }
    }

    return nil
}