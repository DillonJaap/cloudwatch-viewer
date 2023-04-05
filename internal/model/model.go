package model

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	event "clviewer/internal/model/logeventlist"
	vp "clviewer/internal/model/logeventviewport"
	group "clviewer/internal/model/loggrouplist"
)

var (
	modelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder())

	selectedModelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

var (
	groupListSelected = 0
	eventListSelected = 1
	viewportSelected  = 2
)

type Model struct {
	eventsList     event.Model
	viewportEvents vp.Model
	logGroupList   group.Model
	selected       int
}

func InitialModel(ctx context.Context) *Model {
	return &Model{
		eventsList:     event.New(),
		viewportEvents: vp.New("..."),
		logGroupList:   group.New(),
		selected:       eventListSelected,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var (
		logGroupList string
		eventList    string
		message      string
	)

	// TODO clean this up
	switch m.selected {
	case groupListSelected:
		logGroupList = selectedModelStyle.Render(m.logGroupList.View())
		eventList = modelStyle.Render(m.eventsList.View())
		message = modelStyle.Render(m.viewportEvents.View())
	case eventListSelected:
		eventList = selectedModelStyle.Render(m.eventsList.View())
		logGroupList = modelStyle.Render(m.logGroupList.View())
		message = modelStyle.Render(m.viewportEvents.View())
	case viewportSelected:
		message = selectedModelStyle.Render(m.viewportEvents.View())
		eventList = modelStyle.Render(m.eventsList.View())
		logGroupList = modelStyle.Render(m.logGroupList.View())
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			logGroupList,
			eventList,
		),
		message,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.selected = (m.selected + 1) % 3
			return m, nil
		}

		switch m.selected {
		case groupListSelected:
			m.logGroupList, cmd = m.logGroupList.Update(msg)
			return m, cmd
		case eventListSelected:
			m.eventsList, cmd = m.eventsList.Update(msg)
			return m, cmd
		case viewportSelected:
			m.viewportEvents, cmd = m.viewportEvents.Update(msg)
			return m, cmd
		}
	case event.UpdateViewPortContent:
		eventValue := m.eventsList.List.SelectedItem().FilterValue()
		m.viewportEvents.Viewport.SetContent(event.FormatMessage(eventValue, true))
		return m, nil
	}

	m.logGroupList, cmd = m.logGroupList.Update(msg)
	cmds = append(cmds, cmd)

	m.eventsList, cmd = m.eventsList.Update(msg)
	cmds = append(cmds, cmd)

	m.viewportEvents, cmd = m.viewportEvents.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
