package loggrouplist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
)

const listHeight = 14

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// Log group model
var _ tea.Model = Model{}

type Model struct {
	List   list.Model
	Choice string
}

func New() Model {
	itemList := GetLogGroupsAsItemList("/aws/lambda")

	groupList := list.New(itemList, &ItemDelegate{}, 0, 0)
	groupList.SetShowStatusBar(false)
	groupList.SetFilteringEnabled(true)
	groupList.Title = "Log Groups"
	groupList.Styles.PaginationStyle = paginationStyle
	groupList.SetShowHelp(false)

	return Model{
		List:   groupList,
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
		m.List.SetHeight(msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.List.SelectedItem().(Item)
			if ok {
				m.Choice = string(i)
			}

			return m, commands.UpdateEventListItems(m.Choice)
		}
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return "\n" + m.List.View()
}
