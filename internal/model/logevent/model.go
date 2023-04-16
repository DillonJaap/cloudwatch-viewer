package logevent

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
	"clviewer/internal/model/logevent/message"
	"clviewer/internal/model/logevent/timestamp"
)

var (
	doubleBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("08")).
			MarginLeft(1)

	bold = lipgloss.NewStyle().
		Bold(true)

	purpleText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("98"))
)

type Model struct {
	Timestamp      timestamp.Model
	Messages       message.Model
	selectedGroup  string
	selectedStream string
}

// TODO should I add a New() function?

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateWindowSize(msg)
	case tea.KeyMsg:
		return m.handleUpdateKey(msg)
	case commands.UpdateEventListItemsMsg:
		m.selectedGroup = msg.Group
		m.selectedStream = msg.Stream
	case commands.UpdateStreamListItemsMsg:
		m.selectedGroup = msg.Group
	}

	m.Timestamp, cmd = m.Timestamp.Update(msg)
	cmds = append(cmds, cmd)

	m.Messages, cmd = m.Messages.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	logEventView := lipgloss.JoinVertical(
		lipgloss.Left,
		doubleBorder.Render(fmt.Sprintf(
			" %s: %s %s: %s ",
			bold.Render("LogGroup"),
			purpleText.Render(m.selectedGroup),
			bold.Render("LogStream"),
			purpleText.Render(m.selectedStream),
		))+"\n",
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.Timestamp.View(),
			m.Messages.View(),
		),
	)

	return logEventView
}

func (m Model) handleUpdateWindowSize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	const statusBarHeight = 4

	height := msg.Height - statusBarHeight

	timestampWidth := int(float32(msg.Width) / 3.0)
	messageWidth := msg.Width - timestampWidth

	m.Timestamp, cmd = m.Timestamp.Update(tea.WindowSizeMsg{
		Width:  timestampWidth,
		Height: height,
	})
	cmds = append(cmds, cmd)

	m.Messages, cmd = m.Messages.Update(tea.WindowSizeMsg{
		Width:  messageWidth,
		Height: height,
	})
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) handleUpdateKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch keypress := msg.String(); keypress {
	case "H", "J", "K", "L", "enter", "c", "/":
		m.Timestamp, cmd = m.Timestamp.Update(msg)
		return m, cmd
	default:
		m.Messages, cmd = m.Messages.Update(msg)
		return m, cmd
	}
}
