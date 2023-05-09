package logevent

import "github.com/charmbracelet/bubbles/key"

const spacebar = " "

type keyMap struct {
	PrevItem     key.Binding
	NextItem     key.Binding
	ScrollUp     key.Binding
	ScrollDown   key.Binding
	Left         key.Binding
	Right        key.Binding
	Help         key.Binding
	Quit         key.Binding
	Filter       key.Binding
	LoadMore     key.Binding
	Collapse     key.Binding
	CollapseAll  key.Binding
	NextWindow   key.Binding
	PrevWindow   key.Binding
	Copy         key.Binding
	PageDown     key.Binding
	PageUp       key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	Reload       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
		{k.NextWindow, k.PrevWindow, k.PrevItem, k.NextItem},
		{k.ScrollUp, k.ScrollDown, k.PageUp, k.PageDown},
		{k.HalfPageUp, k.HalfPageDown, k.Collapse, k.CollapseAll},
		{k.Filter, k.Copy, k.LoadMore, k.Reload},
	}
}

var keys = keyMap{
	PrevItem: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "next item"),
	),
	NextItem: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "prev item"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("shift+up", "K"),
		key.WithHelp("shift+↑/K", "scroll up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("shift+down", "J"),
		key.WithHelp("shift+↓/J", "scroll down"),
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
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy"),
	),
	Collapse: key.NewBinding(
		key.WithKeys(spacebar),
		key.WithHelp("spacebar", "toggle collapse"),
	),
	CollapseAll: key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "toggle collapse all"),
	),
	NextWindow: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next window"),
	),
	PrevWindow: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev window"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("f/pgdn", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("b/pgup", "page up"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("u", "ctrl+u"),
		key.WithHelp("u", "½ page up"),
	),
	HalfPageDown: key.NewBinding(
		key.WithKeys("d", "ctrl+d"),
		key.WithHelp("d", "½ page down"),
	),
	Reload: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "reload events"),
	),
}
