package logstreamlist

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
	List           list.Model
	SelectedStream string
	currentGroup   string
	padding        int
}

func New(
	title string,
) Model {
	streamList := list.New([]list.Item{}, &ItemDelegate{}, 0, 0)

	streamList.SetShowStatusBar(false)
	streamList.SetFilteringEnabled(true)
	streamList.SetShowHelp(false)

	streamList.Title = title
	streamList.Styles.Title = titleStyle
	streamList.Styles.PaginationStyle = paginationStyle

	return Model{
		List:           streamList,
		SelectedStream: "",
		currentGroup:   "",
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

		log.Printf("%+v\n", m.padding)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.List.SelectedItem().(Item)
			if ok {
				m.SelectedStream = string(i)
			}

			return m, commands.UpdateEventListItems(m.currentGroup, m.SelectedStream)
		}
	case commands.UpdateStreamListItemsMsg:
		m.currentGroup = msg.Group
		cmd = m.UpdateStreamItems(m.currentGroup)
		cmds = append(cmds, cmd)
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		PaddingRight(m.List.Width() - lipgloss.Width(m.List.View())).
		Render(m.List.View())
}

func (m *Model) UpdateStreamItems(groupPattern string) tea.Cmd {
	// reset list
	m.SelectedStream = ""
	m.List.FilterInput.SetCursor(0)
	itemList := GetLogStreamsAsItemList(groupPattern)
	return m.List.SetItems(itemList)
}
