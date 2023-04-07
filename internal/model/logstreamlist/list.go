package logstreamlist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
)

const listHeight = 14

var (
	titleStyle = lipgloss.
			NewStyle().
			Background(lipgloss.Color("98")).
			Foreground(lipgloss.Color("230")).
			PaddingLeft(1).
			PaddingRight(1)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// Log Stream model
var _ tea.Model = Model{}

type Model struct {
	List   list.Model
	Choice string
}

func New(
	title string,
	name string,
) Model {
	streamList := list.New([]list.Item{}, &ItemDelegate{}, 0, 0)

	streamList.SetShowStatusBar(false)
	streamList.SetFilteringEnabled(true)
	streamList.SetShowHelp(false)

	streamList.Title = title
	streamList.Styles.Title = titleStyle
	streamList.Styles.PaginationStyle = paginationStyle

	return Model{
		List:   streamList,
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

			return m, commands.UpdateStreamListItems(m.Choice)
		}
	case commands.UpdateStreamListItemsMsg:
		return m, m.UpdateStreamItems(msg.Group)
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.List.View()
}

func (m *Model) UpdateStreamItems(groupPattern string) tea.Cmd {
	itemList := GetLogStreamsAsItemList(groupPattern)
	return m.List.SetItems(itemList)
}
