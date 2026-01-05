package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Expand       key.Binding
	Collapse     key.Binding
	ToggleSelect key.Binding
	Play         key.Binding
	Quit         key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "向上移动"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "向下移动"),
	),
	Expand: key.NewBinding(
		key.WithKeys("l", "enter"),
		key.WithHelp("l/enter", "展开"),
	),
	Collapse: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "收起"),
	),
	ToggleSelect: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "选择/取消"),
	),
	Play: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "播放选中项"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "退出"),
	),
}
