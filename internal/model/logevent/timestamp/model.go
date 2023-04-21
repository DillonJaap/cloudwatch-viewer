package timestamp

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/cloudwatch/event"
	"clviewer/internal/commands"
)

var (
	titleStyle = lipgloss.
			NewStyle().
			Background(lipgloss.Color("98")).
			Foreground(lipgloss.Color("230")).
			PaddingLeft(1).
			PaddingRight(1)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	selectedItemStyleViewPort = lipgloss.NewStyle().Foreground(lipgloss.Color("127"))
)

type Model struct {
	List           list.Model
	Choice         string
	CollapseAll    bool
	eventPaginator *event.Paginator
}

func (m Model) Init() tea.Cmd {
	return nil
}

type LoadMoreEventsMsg []types.OutputLogEvent

type ResetMsg struct{}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height)
		return m, nil
	case tea.KeyMsg:
		return m.updateKeyMsg(msg)
	case LoadMoreEventsMsg:
		m.List.SetItems(append(
			m.List.Items(),
			logEventsToItemList(msg)...,
		))
	case ResetMsg:
		m.List.Update([]list.Item{})
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.List.View()
}

// updateKeyMsg updates model based on the tea.KeyMsg
func (m Model) updateKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	index := m.List.Index()
	switch keypress := msg.String(); keypress {
	// Toggle Collapse on Item
	case "enter":
		item, ok := m.List.SelectedItem().(Item)
		if ok {
			m.Choice = item.Message
		}
		m.ItemMetaData[index].Collapsed = !m.ItemMetaData[index].Collapsed
		cmd = commands.UpdateViewPort(
			m.getItemListAsStringArray(),
			m.ItemMetaData[index].lineNum,
		)
		return m, cmd
	// Toggle Collapse all
	case "J":
		m.List.CursorDown()
		cmd = commands.UpdateViewPort(
			m.getItemListAsStringArray(),
			m.getLineNumber(),
		)
		return m, cmd
	case "K":
		m.List.CursorUp()
		cmd = commands.UpdateViewPort(
			m.getItemListAsStringArray(),
			m.getLineNumber(),
		)
		return m, cmd
	// all other keystrokes get handled by the list Model
	// and then the viewport gets updated
	default:
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)

		cmd = commands.UpdateViewPort(
			m.getItemListAsStringArray(),
			m.getLineNumber(),
		)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

}

func (m *Model) getItemListAsStringArray() []string {
	var list []string
	height := 0

	for index, item := range m.List.Items() {
		// TODO do I need to format this?
		formattedItem := FormatMessage(
			item.FilterValue(),
			!m.ItemMetaData[index].Collapsed,
		)

		if m.List.Index() == index {
			list = append(list, selectedItemStyleViewPort.Render(formattedItem))
		} else {
			list = append(list, formattedItem)
		}

		height += lipgloss.Height(formattedItem)
		m.ItemMetaData[index].lineNum = height
	}
	return list
}

func (m Model) getLineNumber() int {
	lineNum := 0
	if len(m.ItemMetaData) > m.List.Index()+1 {
		lineNum = m.ItemMetaData[m.List.Index()].lineNum
	}
	return lineNum
}
