package logeventlist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UpdateViewPortContent struct{}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type Model struct {
	List   list.Model
	Choice string
}

func New() Model {
	//defaultWidth := 20
	itemList := GetLogEventsAsItemList()
	itemList = formatList(itemList, false)

	eventList := list.New(itemList, &ItemDelegate{}, 0, 0)
	eventList.SetShowStatusBar(false)
	eventList.SetFilteringEnabled(true)
	eventList.Title = "Log Events"
	eventList.Styles.PaginationStyle = paginationStyle
	eventList.Styles.HelpStyle = helpStyle

	return Model{
		List:   eventList,
		Choice: "",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			e, ok := m.List.SelectedItem().(Item)
			if ok {
				m.Choice = e.Message
			}
			return m, nil
		default:
			cmd = func() tea.Msg { return UpdateViewPortContent{} }
			cmds = append(cmds, cmd)
		}
	}
	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.Choice))
	}
	return "\n" + m.List.View()
}
