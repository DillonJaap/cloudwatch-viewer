package model

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
	event "clviewer/internal/model/logeventlist"
	vp "clviewer/internal/model/logeventviewport"
	group "clviewer/internal/model/loggrouplist"
)

var (
	modelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder())

	selectedModelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("69"))
)

const (
	groupListSelected = 0
	eventListSelected = 1
	viewportSelected  = 2
)

type Model struct {
	eventsList     event.Model
	viewportEvents vp.Model
	logGroupList   group.Model
	keys           keyMap
	help           help.Model
	selected       int
}

func InitialModel(ctx context.Context) *Model {
	return &Model{
		eventsList:     event.New(),
		viewportEvents: vp.New("..."),
		logGroupList:   group.New(),
		keys:           keys,
		help:           help.New(),
		selected:       eventListSelected,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	var (
		logGroupList       string
		eventList          string
		logMessageViewPort string
	)

	helpView := m.help.View(m.keys)

	// TODO clean this up
	switch m.selected {
	case groupListSelected:
		logGroupList = selectedModelStyle.Render(m.logGroupList.View())
		eventList = modelStyle.Render(m.eventsList.View())
		logMessageViewPort = modelStyle.Render(m.viewportEvents.View())
	case eventListSelected:
		eventList = selectedModelStyle.Render(m.eventsList.View())
		logGroupList = modelStyle.Render(m.logGroupList.View())
		logMessageViewPort = modelStyle.Render(m.viewportEvents.View())
	case viewportSelected:
		logMessageViewPort = selectedModelStyle.Render(m.viewportEvents.View())
		eventList = modelStyle.Render(m.eventsList.View())
		logGroupList = modelStyle.Render(m.logGroupList.View())
	}

	logGroupAndEvents := lipgloss.JoinVertical(
		lipgloss.Center,
		logGroupList,
		eventList,
	)

	logView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		logGroupAndEvents,
		logMessageViewPort,
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logView,
		helpView,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.selected = (m.selected + 1) % 3
			return m, nil
		default:
			return m.updateKeyMsg(msg)
		}
	case tea.WindowSizeMsg:
		return m.updateWindowSizes(msg)
	case commands.UpdateViewPortContentMsg:
		return m.viewportEvents.Update(msg)
	default:
		return m.updateSubModules(msg)
	}
}

func (m *Model) updateWindowSizes(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	var (
		cmd   tea.Cmd
		cmds  []tea.Cmd
		model tea.Model
	)

	height := msg.Height - lipgloss.Height(m.help.View(keys))
	width := msg.Width

	borderMarginSize := 4 // subtract 4 for border

	leftHandWidth := int(float32(width / 3))
	rightHandWidth := width - leftHandWidth

	logGroupListHeight := int(float32(height) / 4)
	logEventListHeight := height - logGroupListHeight
	eventViewPortHeight := height

	os.WriteFile(
		"log",
		[]byte(fmt.Sprintf(
			"height: %d, width: %d\n"+
				"groupheight: %d, eventheight: %d\n",
			height,
			width,
			logGroupListHeight,
			logEventListHeight,
		)),
		0665,
	)

	model, cmd = m.logGroupList.Update(tea.WindowSizeMsg{
		Width:  leftHandWidth - borderMarginSize,
		Height: logGroupListHeight - borderMarginSize,
	})
	m.logGroupList = model.(group.Model)
	cmds = append(cmds, cmd)

	model, cmd = m.eventsList.Update(tea.WindowSizeMsg{
		Width:  leftHandWidth - borderMarginSize,
		Height: logEventListHeight - borderMarginSize,
	})
	m.eventsList = model.(event.Model)
	cmds = append(cmds, cmd)

	model, cmd = m.viewportEvents.Update(tea.WindowSizeMsg{
		Width:  rightHandWidth - borderMarginSize,
		Height: eventViewPortHeight - borderMarginSize,
	})
	m.viewportEvents = model.(vp.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateKeyMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	var model tea.Model

	switch m.selected {
	case groupListSelected:
		model, cmd = m.logGroupList.Update(msg)
		m.logGroupList = model.(group.Model)
	case eventListSelected:
		model, cmd = m.eventsList.Update(msg)
		m.eventsList = model.(event.Model)
	case viewportSelected:
		model, cmd = m.viewportEvents.Update(msg)
		m.viewportEvents = model.(vp.Model)
	}
	return m, cmd
}

func (m *Model) updateSubModules(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd   tea.Cmd
		cmds  []tea.Cmd
		model tea.Model
	)

	model, cmd = m.logGroupList.Update(msg)
	m.logGroupList = model.(group.Model)
	cmds = append(cmds, cmd)

	model, cmd = m.eventsList.Update(msg)
	m.eventsList = model.(event.Model)
	cmds = append(cmds, cmd)

	model, cmd = m.viewportEvents.Update(msg)
	m.viewportEvents = model.(vp.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
