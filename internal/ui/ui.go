package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ayazumi/biliCLI/internal/model"
	"github.com/ayazumi/biliCLI/internal/tree"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("237")).
			Foreground(lipgloss.Color("229"))

	cursorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type UI struct {
	model  *model.Model
	keys   KeyMap
	help   help.Model
	width  int
	height int
}

func NewUI(m *model.Model) *UI {
	return &UI{
		model: m,
		keys:  DefaultKeyMap,
		help:  help.New(),
	}
}

func (u *UI) Init() tea.Cmd {
	return nil
}

func (u *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, u.keys.Quit):
			return u, tea.Quit
		case key.Matches(msg, u.keys.Up):
			u.model.MoveCursor(-1)
		case key.Matches(msg, u.keys.Down):
			u.model.MoveCursor(1)
		case key.Matches(msg, u.keys.Expand):
			cursor := u.model.GetCursor()
			node := u.model.GetVisibleNodes()[cursor]
			if !node.IsLeaf() {
				u.model.ToggleExpand(cursor, true)
			} else {
				// 播放叶子节点
				return u, u.playNode(node)
			}
		case key.Matches(msg, u.keys.Collapse):
			cursor := u.model.GetCursor()
			node := u.model.GetVisibleNodes()[cursor]
			if !node.IsLeaf() && node.Expanded {
				u.model.ToggleExpand(cursor, false)
			}
		case key.Matches(msg, u.keys.ToggleSelect):
			cursor := u.model.GetCursor()
			u.model.ToggleSelection(cursor)
		case key.Matches(msg, u.keys.Play):
			return u, u.playSelected()
		}

	case tea.WindowSizeMsg:
		u.width = msg.Width
		u.height = msg.Height
		u.help.Width = msg.Width
	}

	return u, nil
}

func (u *UI) View() string {
	if u.model == nil {
		return "加载中..."
	}

	nodes := u.model.GetVisibleNodes()
	cursor := u.model.GetCursor()
	selected := u.model.GetSelectedCIDs()

	var lines []string
	lines = append(lines, titleStyle.Render("BiliCLI - 树形视频浏览器"))
	lines = append(lines, "")

	// 计算可显示的行数
	maxLines := u.height - 4 // 为标题和帮助信息留空间
	start := 0
	if cursor > maxLines/2 {
		start = cursor - maxLines/2
	}
	if start+maxLines > len(nodes) {
		start = len(nodes) - maxLines
	}
	if start < 0 {
		start = 0
	}

	// 显示节点
	for i := start; i < len(nodes) && i < start+maxLines; i++ {
		node := nodes[i]
		line := u.renderNode(node, i == cursor, selected)
		lines = append(lines, line)
	}

	// 添加帮助信息
	lines = append(lines, "")
	helpView := u.help.View(u.keys)
	lines = append(lines, helpStyle.Render(helpView))

	return strings.Join(lines, "\n")
}

func (u *UI) renderNode(node tree.Node, isCursor bool, selected map[uint64]bool) string {
	line := node.Display()

	if isCursor {
		line = cursorStyle.Render(">") + " " + line
	} else {
		line = "  " + line
	}

	// 如果是叶子节点且被选中，高亮显示
	if node.IsLeaf() && selected[node.CID] {
		line = selectedStyle.Render(line)
	}

	return line
}

func (u *UI) playNode(node tree.Node) tea.Cmd {
	if node.IsLeaf() {
		return func() tea.Msg {
			// 这里执行播放命令
			fmt.Printf("播放视频: %s (CID: %d)\n", node.Name, node.CID)
			return nil
		}
	}
	return nil
}

func (u *UI) playSelected() tea.Cmd {
	selected := u.model.GetSelectedCIDs()
	if len(selected) == 0 {
		return nil
	}

	return func() tea.Msg {
		fmt.Printf("播放选中的视频: %v\n", selected)
		return nil
	}
}
