package model

import (
	"context"
	"math"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
	"clviewer/internal/keymap"
	event "clviewer/internal/model/logevent"
	"clviewer/internal/model/logevent/message"
	"clviewer/internal/model/logevent/timestamp"
	group "clviewer/internal/model/loggrouplist"
	stream "clviewer/internal/model/logstreamlist"
)

var (
	modelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("08"))

	selectedModelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("69"))

	tuiBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("08"))

	doubleBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("08")).
			MarginLeft(1)

	bold = lipgloss.NewStyle().
		Bold(true)

	purpleText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("98"))
)

const (
	groupListSelected  = 0
	streamListSelected = 1
	eventListSelected  = 2
	numWindows         = 3
)

type Model struct {
	logEvent  event.Model
	logGroup  group.Model
	logStream stream.Model
	help      help.Model
	keyMaps   []help.KeyMap
	selected  int
}

func New(ctx context.Context) *Model {
	helpModel := help.New()
	helpModel.ShowAll = true
	return &Model{
		logEvent: event.Model{
			Timestamp: timestamp.New(
				"Timestamps",
				false,
			),
			Messages: message.New(
				"Log Messages",
				"...",
			),
		},
		logGroup: group.New(
			"Log Groups",
			"/aws/lambda/dev-djaap",
		),
		logStream: stream.New(
			"Log Streams",
		),
		keyMaps: []help.KeyMap{
			keys,
			keys,
			keymap.Keys,
		},
		help:     helpModel,
		selected: eventListSelected,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	logGroupList := m.logGroup.View()
	logStreamList := m.logStream.View()
	logEventView := m.logEvent.View()

	helpView := m.help.View(m.keyMaps[m.selected])

	switch m.selected {
	case groupListSelected:
		logGroupList = selectedModelStyle.Render(logGroupList)
		logEventView = modelStyle.Render(logEventView)
		logStreamList = modelStyle.Render(logStreamList)
	case streamListSelected:
		logStreamList = selectedModelStyle.Render(logStreamList)
		logEventView = modelStyle.Render(logEventView)
		logGroupList = modelStyle.Render(logGroupList)
	case eventListSelected:
		logEventView = selectedModelStyle.Render(logEventView)
		logGroupList = modelStyle.Render(logGroupList)
		logStreamList = modelStyle.Render(logStreamList)
	}

	logLists := lipgloss.JoinVertical(
		lipgloss.Left,
		logGroupList,
		logStreamList,
	)

	logListsAndEvents := lipgloss.JoinHorizontal(
		lipgloss.Left,
		logLists,
		logEventView,
	)

	return tuiBorder.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			helpView,
			logListsAndEvents,
		),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "?":
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case "tab":
			m.selected = (m.selected + 1) % numWindows
			return m, nil
		case "shift+tab":
			m.selected = int(math.Abs(float64((m.selected - 1) % numWindows)))
			return m, nil
		default:
			return m.updateKeyMsg(msg)
		}
	case tea.WindowSizeMsg:
		return m.updateWindowSizes(msg)
	case commands.UpdateViewPortContentMsg:
		m.logEvent.Update(msg)
	}
	return m.updateSubModules(msg)
}

func (m *Model) updateWindowSizes(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	const borderMarginSize = 4 // subtract 4 for border
	const tuiBorder = 2

	height := msg.Height - lipgloss.Height(m.help.View(keys)) - tuiBorder
	width := msg.Width - tuiBorder

	logGroupListWidth := int(float32(width) / 4.0)
	logGroupListHeight := int(float32(height) / 2.0)
	m.logGroup, cmd = m.logGroup.Update(tea.WindowSizeMsg{
		Width:  logGroupListWidth,
		Height: logGroupListHeight - borderMarginSize,
	})
	cmds = append(cmds, cmd)

	logStreamListWidth := int(float32(width) / 4.0)
	logStreamListHeight := height - logGroupListHeight
	m.logStream, cmd = m.logStream.Update(tea.WindowSizeMsg{
		Width:  logStreamListWidth,
		Height: logStreamListHeight - borderMarginSize,
	})
	cmds = append(cmds, cmd)

	eventsWidth := width - logGroupListWidth
	logEventHeight := height
	m.logEvent, cmd = m.logEvent.Update(tea.WindowSizeMsg{
		Width:  eventsWidth - borderMarginSize,
		Height: logEventHeight - borderMarginSize,
	})
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateKeyMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch m.selected {
	case groupListSelected:
		m.logGroup, cmd = m.logGroup.Update(msg)
	case streamListSelected:
		m.logStream, cmd = m.logStream.Update(msg)
	case eventListSelected:
		m.logEvent, cmd = m.logEvent.Update(msg)
	}
	return m, cmd
}

func (m *Model) updateSubModules(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.logGroup, cmd = m.logGroup.Update(msg)
	cmds = append(cmds, cmd)

	m.logStream, cmd = m.logStream.Update(msg)
	cmds = append(cmds, cmd)

	m.logEvent, cmd = m.logEvent.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
