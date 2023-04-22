package timestamp

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/cloudwatch/event"
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

type NextEventMsg struct{}

type PrevEventMsg struct{}

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
	case NextEventMsg:
		m.List.CursorDown()
	case PrevEventMsg:
		m.List.CursorUp()
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
	switch keypress := msg.String(); keypress {
	// all other keystrokes get handled by the list Model
	default:
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
}
