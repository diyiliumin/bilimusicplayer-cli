package tree

import "fmt"

type NodeType int

const (
	NodeGroup NodeType = iota
	NodeTitle
	NodeTab
	NodeItem // 叶子节点
)

type Node struct {
	Type      NodeType
	Depth     int
	Name      string
	CID       uint64 // 仅用于Item节点
	P         uint32 // 用于排序
	Expanded  bool
	ParentIdx int // 用于快速收起（可选）
}

func (n Node) Display() string {
	indent := ""
	for i := 0; i < n.Depth; i++ {
		indent += "  "
	}
	marker := " "
	if n.Type != NodeItem {
		if n.Expanded {
			marker = "▼"
		} else {
			marker = "▶"
		}
	}
	return fmt.Sprintf("%s%s %s", indent, marker, n.Name)
}

func (n Node) IsLeaf() bool {
	return n.Type == NodeItem
}
