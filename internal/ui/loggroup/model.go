package loggroup

import (
	"log"

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
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height)
		return m, nil
	case tea.KeyMsg:
		if isRedrawKey(msg) {
			cmds = append(cmds, commands.RedrawWindows())
		}

		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			// don't apply item if currently setting filter
			if m.List.SettingFilter() {
				m.List, cmd = m.List.Update(msg)
				return m, cmd
			}

			// TODO update to not use key press, but checking to see if
			// selected item changed to send the updateStreamListCommand?
			// Could this logic be moved up to ui/model?
			i, ok := m.List.SelectedItem().(Item)
			if ok {
				m.SelectedGroup = string(i)
			}
			log.Printf("loggroup: %+v", m.SelectedGroup)

			cmds := append(cmds, commands.UpdateStreamListItems(m.SelectedGroup))
			return m, tea.Batch(cmds...)
		}
	}
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

// isRedrawKey checks to see if the keypress should trigger a redraw of the ui
// TODO update to use proper keymap
func isRedrawKey(key tea.KeyMsg) bool {
	if key.String() == "enter" {
		return true
	}
	return false
}
