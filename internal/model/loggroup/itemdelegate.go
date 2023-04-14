package loggroup

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Item Delegate
type ItemDelegate struct{}

func (d ItemDelegate) Height() int { return 1 }

func (d ItemDelegate) Spacing() int { return 0 }

func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", item.getTruncatedDescription(m.Width()-10))

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}
