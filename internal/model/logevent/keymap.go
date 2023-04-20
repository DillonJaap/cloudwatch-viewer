package logevent

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	Help        key.Binding
	Quit        key.Binding
	Filter      key.Binding
	LoadMore    key.Binding
	Collapse    key.Binding
	CollapseAll key.Binding
	NextWindow  key.Binding
	PrevWindow  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Left, k.Right},
		{k.Collapse, k.CollapseAll},
		{k.Filter, k.LoadMore},
		{k.NextWindow, k.PrevWindow},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "next item"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "prev item"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "scroll left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "scroll right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	LoadMore: key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "load more events"),
	),
	Collapse: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "toggle collapse"),
	),
	CollapseAll: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "toggle collapse all"),
	),
	NextWindow: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next window"),
	),
	PrevWindow: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev window"),
	),
}
