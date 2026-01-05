package model

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ayazumi/biliCLI/internal/tree"
)

type Model struct {
	fullTree     []GroupNode
	visibleNodes []tree.Node
	cursor       int
	selected     map[uint64]bool
	searchQuery  string
}

func NewModel() *Model {
	return &Model{
		selected: make(map[uint64]bool),
	}
}

func (m *Model) LoadTreeData(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	var treeData TreeData
	if err := json.Unmarshal(data, &treeData); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	// 初始化所有节点为折叠状态
	for i := range treeData.Groups {
		treeData.Groups[i].Open = false
		for j := range treeData.Groups[i].Titles {
			treeData.Groups[i].Titles[j].Open = false
		}
	}

	m.fullTree = treeData.Groups
	m.rebuildVisible()
	return nil
}

func (m *Model) rebuildVisible() {
	m.visibleNodes = m.buildVisibleFrom(m.fullTree, 0, []bool{})
}

func (m *Model) buildVisibleFrom(groups []GroupNode, depth int, parentPath []bool) []tree.Node {
	var nodes []tree.Node

	for _, g := range groups {
		// 添加组节点
		groupNode := tree.Node{
			Type:     tree.NodeGroup,
			Depth:    depth,
			Name:     g.Name,
			Expanded: g.Open,
		}
		nodes = append(nodes, groupNode)

		// 如果组展开，添加子节点
		if g.Open {
			for _, t := range g.Titles {
				pVal := uint32(0)
				if t.P != nil {
					pVal = *t.P
				}

				titleNode := tree.Node{
					Type:     tree.NodeTitle,
					Depth:    depth + 1,
					Name:     t.Name,
					P:        pVal,
					Expanded: t.Open,
				}
				nodes = append(nodes, titleNode)

				// 如果标题展开，添加标签和项目
				if t.Open {
					for _, tab := range t.Tabs {
						tabNode := tree.Node{
							Type:  tree.NodeTab,
							Depth: depth + 2,
							Name:  tab.Name,
						}
						nodes = append(nodes, tabNode)

						// 添加项目节点（叶子节点）
						for _, item := range tab.Items {
							itemNode := tree.Node{
								Type:  tree.NodeItem,
								Depth: depth + 3,
								Name:  item.Title,
								CID:   item.CID,
								P:     item.P,
							}
							nodes = append(nodes, itemNode)
						}
					}
				}
			}
		}
	}

	return nodes
}

func (m *Model) ToggleExpand(cursor int, expand bool) {
	if cursor < 0 || cursor >= len(m.visibleNodes) {
		return
	}

	node := &m.visibleNodes[cursor]

	// 找到对应的fullTree节点并切换状态
	m.toggleFullTreeNode(node, expand)
	m.rebuildVisible()
}

func (m *Model) toggleFullTreeNode(targetNode *tree.Node, expand bool) {
	// 根据节点类型和名称找到对应的fullTree节点
	// 这里需要遍历fullTree来找到匹配的节点

	if targetNode.Type == tree.NodeGroup {
		for i := range m.fullTree {
			if m.fullTree[i].Name == targetNode.Name {
				m.fullTree[i].Open = expand
				return
			}
		}
	} else if targetNode.Type == tree.NodeTitle {
		// 需要找到对应的组和标题
		for i := range m.fullTree {
			for j := range m.fullTree[i].Titles {
				if m.fullTree[i].Titles[j].Name == targetNode.Name {
					m.fullTree[i].Titles[j].Open = expand
					return
				}
			}
		}
	}
}

func (m *Model) GetVisibleNodes() []tree.Node {
	return m.visibleNodes
}

func (m *Model) GetCursor() int {
	return m.cursor
}

func (m *Model) SetCursor(cursor int) {
	if cursor >= 0 && cursor < len(m.visibleNodes) {
		m.cursor = cursor
	}
}

func (m *Model) MoveCursor(delta int) {
	newCursor := m.cursor + delta
	if newCursor >= 0 && newCursor < len(m.visibleNodes) {
		m.cursor = newCursor
	}
}

func (m *Model) ToggleSelection(cursor int) {
	if cursor < 0 || cursor >= len(m.visibleNodes) {
		return
	}

	node := m.visibleNodes[cursor]
	if node.IsLeaf() {
		if m.selected[node.CID] {
			delete(m.selected, node.CID)
		} else {
			m.selected[node.CID] = true
		}
	}
}

func (m *Model) GetSelectedCIDs() []uint64 {
	var cids []uint64
	for cid := range m.selected {
		cids = append(cids, cid)
	}
	return cids
}
