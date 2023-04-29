package loggroup

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

type Model struct {
	List          list.Model
	SelectedGroup string
	padding       int
}

func New(
	title string,
	groupPattern string,
	intialGroup string,
) Model {
	itemList := GetLogGroupsAsItemList(groupPattern)

	groupList := list.New(itemList, &ItemDelegate{}, 0, 0)

	groupList.SetShowStatusBar(false)
	groupList.SetFilteringEnabled(true)
	groupList.SetShowHelp(false)

	groupList.Title = title
	groupList.Styles.Title = titleStyle
	groupList.Styles.PaginationStyle = paginationStyle

	return Model{
		List:          groupList,
		SelectedGroup: "initialGroup",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
				m.SelectedGroup = string(i)
			}

			return m, commands.UpdateStreamListItems(m.SelectedGroup)
		}
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		PaddingRight(m.List.Width() - lipgloss.Width(m.List.View())).
		Render(m.List.View())
}

func (m Model) HelpView() string {
	return m.List.Styles.HelpStyle.Render(m.List.Help.View(m.List))
}
