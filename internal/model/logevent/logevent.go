package logevent

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
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

func DefaultModel() Model {
	//defaultWidth := 20
	eventList := list.New(
		GetLogEventsAsItemList(),
		&eventDelegate{},
		0,
		0,
	)
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			e, ok := m.List.SelectedItem().(Event)
			if ok {
				m.Choice = e.Message
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.Choice))
	}
	return "\n" + m.List.View()
}
