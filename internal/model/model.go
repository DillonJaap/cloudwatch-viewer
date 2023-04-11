package model

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
	"clviewer/internal/keymap"
	event "clviewer/internal/model/logeventlist"
	vp "clviewer/internal/model/logeventviewport"
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
)

const (
	groupListSelected  = 0
	streamListSelected = 1
	eventListSelected  = 2
	numWindows         = 3
)

type Model struct {
	eventsList     event.Model
	viewportEvents vp.Model
	logGroupList   group.Model
	logStreamList  stream.Model
	help           help.Model
	keyMaps        []help.KeyMap
	selected       int
}

func New(ctx context.Context) *Model {
	helpModel := help.New()
	helpModel.ShowAll = true
	return &Model{
		eventsList: event.New(
			"Timestamps",
			false,
		),
		viewportEvents: vp.New(
			"Log Messages",
			"...",
		),
		logGroupList: group.New(
			"Log Groups",
			"/aws/lambda",
		),
		logStreamList: stream.New(
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
	var logGroupList string
	var logStreamList string

	helpView := m.help.View(m.keyMaps[m.selected])

	logEventView := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf(
			"LogGroup: %s | LogStream: %s\n",
			m.logGroupList.SelectedGroup,
			m.logStreamList.SelectedStream,
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.eventsList.View(),
			m.viewportEvents.View(),
		),
	)

	// TODO clean this up
	switch m.selected {
	case groupListSelected:
		logEventView = modelStyle.Render(logEventView)
		logGroupList = selectedModelStyle.Render(m.logGroupList.View())
		logStreamList = modelStyle.Render(m.logStreamList.View())
	case streamListSelected:
		logEventView = modelStyle.Render(logEventView)
		logGroupList = modelStyle.Render(m.logGroupList.View())
		logStreamList = selectedModelStyle.Render(m.logStreamList.View())
	case eventListSelected:
		logEventView = selectedModelStyle.Render(logEventView)
		logGroupList = modelStyle.Render(m.logGroupList.View())
		logStreamList = modelStyle.Render(m.logStreamList.View())
	}

	logLists := lipgloss.JoinHorizontal(
		lipgloss.Center,
		logGroupList,
		logStreamList,
	)

	logGroupAndEvents := lipgloss.JoinVertical(
		lipgloss.Left,
		logLists,
		logEventView,
	)

	return tuiBorder.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			helpView,
			logGroupAndEvents,
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
		m.viewportEvents.Update(msg)
	}
	log.Printf("%+v\n", m.viewportEvents)
	return m.updateSubModules(msg)
}

func (m *Model) updateWindowSizes(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	height := msg.Height - lipgloss.Height(m.help.View(keys)) - 2 // TODO add const for 2
	width := msg.Width

	borderMarginSize := 4 // subtract 4 for border

	logListHeight := int(float32(height) / 4.0)
	logEventListHeight := height - logListHeight - 1  // TODO add const
	eventViewPortHeight := height - logListHeight - 1 // TODO add const

	logGroupListWidth := int(float32(width) / 2.0)
	logStreamListWidth := width - logGroupListWidth

	eventsLeftHandWidth := int(float32(width) / 3.0)
	eventsRightHandWidth := width - eventsLeftHandWidth

	m.logGroupList, cmd = m.logGroupList.Update(tea.WindowSizeMsg{
		Width:  logGroupListWidth,
		Height: logListHeight - borderMarginSize,
	})
	cmds = append(cmds, cmd)

	m.logStreamList, cmd = m.logStreamList.Update(tea.WindowSizeMsg{
		Width:  logStreamListWidth,
		Height: logListHeight - borderMarginSize,
	})
	cmds = append(cmds, cmd)

	m.eventsList, cmd = m.eventsList.Update(tea.WindowSizeMsg{
		Width:  eventsLeftHandWidth - borderMarginSize,
		Height: logEventListHeight - borderMarginSize,
	})
	cmds = append(cmds, cmd)

	m.viewportEvents, cmd = m.viewportEvents.Update(tea.WindowSizeMsg{
		Width:  eventsRightHandWidth - borderMarginSize,
		Height: eventViewPortHeight - borderMarginSize,
	})
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateKeyMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch m.selected {
	case groupListSelected:
		m.logGroupList, cmd = m.logGroupList.Update(msg)
	case streamListSelected:
		m.logStreamList, cmd = m.logStreamList.Update(msg)
	case eventListSelected:
		m.eventsList, cmd = m.eventsList.Update(msg)
	}
	return m, cmd
}

func (m *Model) updateSubModules(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.logGroupList, cmd = m.logGroupList.Update(msg)
	cmds = append(cmds, cmd)

	m.logStreamList, cmd = m.logStreamList.Update(msg)
	cmds = append(cmds, cmd)

	m.eventsList, cmd = m.eventsList.Update(msg)
	cmds = append(cmds, cmd)

	m.viewportEvents, cmd = m.viewportEvents.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
