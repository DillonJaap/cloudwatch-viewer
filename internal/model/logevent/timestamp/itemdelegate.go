package timestamp

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemDelegate struct{}

func (i *ItemDelegate) Height() int { return 1 }

func (i *ItemDelegate) Spacing() int { return 0 }

func (i *ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	m.VisibleItems()
	return nil
}

func (i *ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var str string

	item, ok := listItem.(Item)
	if ok {
		str = fmt.Sprintf("%s", item.getTruncatedTimeStamp(m.Width()-10))
	} else {
		str = fmt.Sprintf("%s", listItem.FilterValue())
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}
