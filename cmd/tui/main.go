package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

const TreeJSONPath = "buildtree/tree.json"

type buildFinishedMsg struct{ err error }

// ========== 数据结构 ==========
type Item struct {
	Title string `json:"title"`
	CID   uint64 `json:"cid"`
}

type TitleNode struct {
	Name  string `json:"name"`
	P     *uint32
	Items []Item
	Open  bool
}

type GroupNode struct {
	Name   string      `json:"name"`
	Titles []TitleNode `json:"titles"`
	Open   bool
}

type NodeType int

const (
	NodeGroup NodeType = iota
	NodeTitle
	NodeItem
)

type TreeNode struct {
	Type     NodeType
	Depth    int
	Name     string
	CID      uint64
	Expanded bool
	groupIdx, titleIdx int
	items    []Item
}

func (n TreeNode) Display() string {
	indent := strings.Repeat("  ", n.Depth)
	marker := " "
	if n.Type == NodeGroup || n.Type == NodeTitle {
		if n.Expanded {
			marker = "▼"
		} else {
			marker = "▶"
		}
	}
	return fmt.Sprintf("%s%s %s", indent, marker, n.Name)
}

// ========== 加载 tree.json ==========
func loadTree() []GroupNode {
	data, err := os.ReadFile(TreeJSONPath)
	if err != nil {
		log.Fatal("❌ 无法读取 tree.json")
	}

	var rawGroups []struct {
		Name   string `json:"name"`
		Titles []struct {
			Name  string `json:"name"`
			P     *uint32
			Tabs  []struct {
				Name  string `json:"name"`
				Items []struct {
					CID uint64 `json:"cid"`
				} `json:"items"`
			} `json:"tabs"`
		} `json:"titles"`
	}
	if err := json.Unmarshal(data, &rawGroups); err != nil {
		log.Fatal("❌ 解析 tree.json 失败:", err)
	}

	var groups []GroupNode
	for _, rg := range rawGroups {
		var titles []TitleNode
		for _, rt := range rg.Titles {
			var items []Item
			for _, tab := range rt.Tabs {
				if len(tab.Items) > 0 {
					items = append(items, Item{
						Title: tab.Name,
						CID:   tab.Items[0].CID,
					})
				}
			}
			titles = append(titles, TitleNode{
				Name:  rt.Name,
				P:     rt.P,
				Items: items,
				Open:  false,
			})
		}
		groups = append(groups, GroupNode{
			Name:   rg.Name,
			Titles: titles,
			Open:   false,
		})
	}
	return groups
}

// ========== 异步构建命令 ==========
func buildTreeCmd() tea.Cmd {
	script := `cd buildtree && ([ -f target/release/buildtree ] && ./target/release/buildtree || cargo run --release)`
	cmd := exec.Command("sh", "-c", script)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return buildFinishedMsg{err}
	})
}

// ========== 播放模式 ==========
type PlayMode int

const (
	PlayModeSequential PlayMode = iota
	PlayModeShuffle
)

func (p PlayMode) String() string {
	switch p {
	case PlayModeSequential:
		return "顺序"
	case PlayModeShuffle:
		return "随机"
	default:
		return "未知"
	}
}

// ========== 播放逻辑 ==========
func playCIDs(cids []uint64, mode PlayMode) tea.Cmd {
	return func() tea.Msg {
		list := make([]uint64, len(cids))
		copy(list, cids)
		if mode == PlayModeShuffle {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
		}
		for _, cid := range list {
			cmd := exec.Command("./play", fmt.Sprintf("%d", cid))
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
		return nil
	}
}

// ========== 状态 ==========
type state int

const (
	StateBuildPrompt state = iota
	StateBuilding
	StateTUI
	StateSearchInput
)

// ========== Model ==========
type model struct {
	state        state
	groups       []GroupNode
	visibleNodes []TreeNode
	allNodes     []TreeNode
	cursor       int
	viewport     viewport.Model
	playMode     PlayMode
	buildError   error

	// 搜索相关
	searchInput  textinput.Model // ← 使用 textinput
	lastSearch   string
	lastMatchIdx int
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "输入关键词..."
	ti.Focus()

	m := model{
		playMode:     PlayModeSequential,
		lastMatchIdx: -1,
		searchInput:  ti,
	}
	if _, err := os.Stat(TreeJSONPath); err == nil {
		m.state = StateTUI
		m.groups = loadTree()
		m.rebuildAllNodes()
		m.rebuildVisible()
		m.initViewport()
	} else {
		m.state = StateBuildPrompt
	}
	return m
}

func (m *model) initViewport() {
	v := viewport.New(80, 10)
	m.viewport = v
}

func (m *model) rebuildAllNodes() {
	var nodes []TreeNode
	for gi, g := range m.groups {
		var groupItems []Item
		for _, t := range g.Titles {
			groupItems = append(groupItems, t.Items...)
		}
		nodes = append(nodes, TreeNode{
			Type:     NodeGroup,
			Depth:    0,
			Name:     g.Name,
			groupIdx: gi,
			items:    groupItems,
		})
		for ti, t := range g.Titles {
			nodes = append(nodes, TreeNode{
				Type:     NodeTitle,
				Depth:    1,
				Name:     t.Name,
				groupIdx: gi,
				titleIdx: ti,
				items:    t.Items,
			})
			for _, item := range t.Items {
				nodes = append(nodes, TreeNode{
					Type:     NodeItem,
					Depth:    2,
					Name:     item.Title,
					CID:      item.CID,
					groupIdx: gi,
					titleIdx: ti,
				})
			}
		}
	}
	m.allNodes = nodes
}

func (m *model) rebuildVisible() {
	var nodes []TreeNode
	for gi, g := range m.groups {
		var groupItems []Item
		for _, t := range g.Titles {
			groupItems = append(groupItems, t.Items...)
		}
		nodes = append(nodes, TreeNode{
			Type:     NodeGroup,
			Depth:    0,
			Name:     g.Name,
			Expanded: g.Open,
			groupIdx: gi,
			items:    groupItems,
		})
		if !g.Open {
			continue
		}
		for ti, t := range g.Titles {
			nodes = append(nodes, TreeNode{
				Type:     NodeTitle,
				Depth:    1,
				Name:     t.Name,
				Expanded: t.Open,
				groupIdx: gi,
				titleIdx: ti,
				items:    t.Items,
			})
			if !t.Open {
				continue
			}
			for _, item := range t.Items {
				nodes = append(nodes, TreeNode{
					Type:     NodeItem,
					Depth:    2,
					Name:     item.Title,
					CID:      item.CID,
					groupIdx: gi,
					titleIdx: ti,
				})
			}
		}
	}
	m.visibleNodes = nodes
	if m.state == StateTUI || m.state == StateSearchInput {
		m.refreshViewport()
	}
}

func (m *model) refreshViewport() {
	content := m.renderContent()
	m.viewport.SetContent(content)
	m.adjustViewport()
}

func (m *model) renderContent() string {
	if len(m.visibleNodes) == 0 {
		return "⚠️ 无数据"
	}
	var lines []string
	for i, node := range m.visibleNodes {
		line := node.Display()
		if i == m.cursor {
			line = "> " + line
		} else {
			line = "  " + line
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (m *model) adjustViewport() {
	height := m.viewport.Height
	if height <= 0 || len(m.visibleNodes) == 0 {
		return
	}
	top := m.viewport.YOffset
	bottom := top + int(height) - 1
	if m.cursor < top {
		m.viewport.SetYOffset(m.cursor)
	} else if m.cursor > bottom {
		m.viewport.SetYOffset(m.cursor - int(height) + 1)
	}
}

func (m *model) jumpToAllNode(allIdx int) {
	node := m.allNodes[allIdx]

	m.groups[node.groupIdx].Open = true
	if node.Type == NodeItem || node.Type == NodeTitle {
		m.groups[node.groupIdx].Titles[node.titleIdx].Open = true
	}

	m.rebuildVisible()

	for j, vis := range m.visibleNodes {
		if vis.Type == node.Type &&
			vis.groupIdx == node.groupIdx &&
			vis.Name == node.Name &&
			(vis.Type != NodeItem || vis.titleIdx == node.titleIdx) {
			m.cursor = j
			m.refreshViewport()
			return
		}
	}
	m.cursor = 0
	m.refreshViewport()
}

func (m *model) searchAndJump(query string) {
	if query == "" {
		m.lastSearch = ""
		m.lastMatchIdx = -1
		return
	}

	lowerQuery := strings.ToLower(query)
	found := false
	for i, node := range m.allNodes {
		if strings.Contains(strings.ToLower(node.Name), lowerQuery) {
			m.lastSearch = query
			m.lastMatchIdx = i
			found = true
			m.jumpToAllNode(i)
			break
		}
	}

	if !found {
		m.lastSearch = ""
		m.lastMatchIdx = -1
	}
}

func (m model) helpView() string {
	modeStr := m.playMode.String()
	return fmt.Sprintf("\n\nh=收起  l=展开  j/k=上下  Enter=播放  m=切换模式(%s)  q=退出  b=同步列表  /=搜索（n=next）", modeStr)
}

// ========== Bubble Tea ==========
func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case buildFinishedMsg:
		if msg.err != nil {
			m.buildError = msg.err
			m.state = StateBuildPrompt
			return m, nil
		}
		_ = syscall.Exec(os.Args[0], os.Args, os.Environ())
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case StateBuildPrompt:
			switch msg.String() {
			case "b", "B":
				m.state = StateTUI 
				return m, buildTreeCmd()
			case "q", "ctrl+c":
				return m, tea.Quit
			}

		case StateSearchInput:
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			switch key := msg.String(); key {
			case "enter":
				query := m.searchInput.Value()
				m.searchAndJump(query)
				m.state = StateTUI
				m.searchInput.Blur()
				return m, nil
			case "esc":
				m.state = StateTUI
				m.searchInput.Blur()
				return m, nil
			}
			return m, cmd // ← 必须返回 cmd！

		case StateTUI:
			switch key := msg.String(); key {
			case "b", "B":
				m.state = StateTUI
				return m, buildTreeCmd()
			case "q", "ctrl+c":
				exec.Command("pkill", "-f", "play").Run()
				return m, tea.Quit
			case "j":
				if m.cursor < len(m.visibleNodes)-1 {
					m.cursor++
					m.refreshViewport()
				}
			case "k":
				if m.cursor > 0 {
					m.cursor--
					m.refreshViewport()
				}
			case "l":
				node := m.visibleNodes[m.cursor]
				if node.Type == NodeGroup {
					m.groups[node.groupIdx].Open = true
					m.rebuildVisible()
				} else if node.Type == NodeTitle {
					m.groups[node.groupIdx].Titles[node.titleIdx].Open = true
					m.rebuildVisible()
				}
			case "h":
				if len(m.visibleNodes) == 0 {
					break
				}
				curr := m.visibleNodes[m.cursor]

				var targetGroupIdx, targetTitleIdx *int
				var targetType NodeType
				foundTarget := false

				switch curr.Type {
				case NodeGroup:
					if m.groups[curr.groupIdx].Open {
						m.groups[curr.groupIdx].Open = false
						targetType = NodeGroup
						gi := curr.groupIdx
						targetGroupIdx = &gi
						foundTarget = true
					}
				case NodeTitle:
					if m.groups[curr.groupIdx].Titles[curr.titleIdx].Open {
						m.groups[curr.groupIdx].Titles[curr.titleIdx].Open = false
						targetType = NodeTitle
						gi, ti := curr.groupIdx, curr.titleIdx
						targetGroupIdx = &gi
						targetTitleIdx = &ti
						foundTarget = true
					} else {
						m.groups[curr.groupIdx].Open = false
						targetType = NodeGroup
						gi := curr.groupIdx
						targetGroupIdx = &gi
						foundTarget = true
					}
				case NodeItem:
					title := &m.groups[curr.groupIdx].Titles[curr.titleIdx]
					group := &m.groups[curr.groupIdx]
					if title.Open {
						title.Open = false
						targetType = NodeTitle
						gi, ti := curr.groupIdx, curr.titleIdx
						targetGroupIdx = &gi
						targetTitleIdx = &ti
						foundTarget = true
					} else if group.Open {
						group.Open = false
						targetType = NodeGroup
						gi := curr.groupIdx
						targetGroupIdx = &gi
						foundTarget = true
					}
				}

				if foundTarget {
					m.rebuildVisible()
					newCursor := -1
					for i, node := range m.visibleNodes {
						if node.Type == targetType && node.groupIdx == *targetGroupIdx {
							if targetType == NodeGroup {
								newCursor = i
								break
							} else if targetType == NodeTitle && node.titleIdx == *targetTitleIdx {
								newCursor = i
								break
							}
						}
					}
					if newCursor != -1 {
						m.cursor = newCursor
						m.refreshViewport()
					} else {
						m.cursor = 0
						m.refreshViewport()
					}
				}

			case "enter":
				node := m.visibleNodes[m.cursor]
				var cids []uint64
				if node.Type == NodeItem {
					cids = []uint64{node.CID}
				} else {
					for _, item := range node.items {
						cids = append(cids, item.CID)
					}
				}
				m.lastSearch = ""
				m.lastMatchIdx = -1
				if len(cids) > 0 {
					return m, playCIDs(cids, m.playMode)
				}

			case "m":
				if m.playMode == PlayModeSequential {
					m.playMode = PlayModeShuffle
				} else {
					m.playMode = PlayModeSequential
				}
				m.refreshViewport()

			case "/":
				m.state = StateSearchInput
				m.searchInput.SetValue("")
				m.searchInput.Focus()
				return m, nil

			case "n":
				if m.lastSearch == "" || len(m.allNodes) == 0 {
					break
				}
				lowerQuery := strings.ToLower(m.lastSearch)
				start := m.lastMatchIdx + 1
				found := false
				for i := start; i < len(m.allNodes); i++ {
					if strings.Contains(strings.ToLower(m.allNodes[i].Name), lowerQuery) {
						m.lastMatchIdx = i
						found = true
						break
					}
				}
				if !found {
					for i := 0; i < start; i++ {
						if strings.Contains(strings.ToLower(m.allNodes[i].Name), lowerQuery) {
							m.lastMatchIdx = i
							found = true
							break
						}
					}
				}
				if found {
					m.jumpToAllNode(m.lastMatchIdx)
				}
			}
		}

	case error:
		if m.state == StateBuilding {
			m.buildError = msg
			m.state = StateBuildPrompt
		}
		return m, nil

	case tea.WindowSizeMsg:
		if m.state == StateTUI || m.state == StateSearchInput {
			helpHeight := 3
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - helpHeight
			m.refreshViewport()
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case StateBuildPrompt:
		msg := "未检测到 buildtree/tree.json\n\n按 B 构建，Q 退出"
		if m.buildError != nil {
			msg += "\n\n❗ " + m.buildError.Error()
		}
		return msg
	case StateBuilding:
		return "正在同步列表...\n"
	case StateSearchInput:
		return "\n搜索: " + m.searchInput.View() + "\n\n（按 Enter 搜索，Esc 取消）"
	case StateTUI:
		return m.viewport.View() + m.helpView()
	default:
		return "未知状态"
	}
}

func main() {
	defer exec.Command("pkill", "-f", "play").Run()
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
