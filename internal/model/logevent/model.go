package logevent

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/cloudwatch/event"
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
	eventPaginator *event.Paginator
	selectedGroup  string
	selectedStream string
	selectedEvent  int
	help           help.Model
}

func New(timestampe timestamp.Model, msg message.Model, help help.Model) Model {
	return Model{
		Timestamp:      timestampe,
		Messages:       message.Model{},
		eventPaginator: nil,
		selectedGroup:  "",
		selectedStream: "",
		selectedEvent:  0,
		help:           help,
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
		return m.handleUpdateWindowSize(msg)
	case tea.KeyMsg:
		return m.handleUpdateKey(msg)
	case commands.UpdateStreamListItemsMsg:
		m.selectedGroup = msg.Group
		return m, nil
	case commands.UpdateEventListItemsMsg:
		m.selectedGroup = msg.Group
		m.selectedStream = msg.Stream

		// get a new paginator for our log group & stream
		paginator := event.New(
			context.Background(),
			m.selectedGroup,
			m.selectedStream,
		)
		m.eventPaginator = &paginator

		{ // reset data
			m.Timestamp, cmd = m.Timestamp.Update(timestamp.ResetMsg{})
			cmds = append(cmds, cmd)
			m.Messages, cmd = m.Messages.Update(message.ResetMsg{})
			cmds = append(cmds, cmd)
		}

		// get initial set of events
		cmd = m.loadMoreEvents()
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
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

	helpView := m.help.View(keys)

	logEventView = lipgloss.JoinVertical(
		lipgloss.Center,
		logEventView,
		helpView,
	)

	return logEventView
}

func (m Model) handleUpdateWindowSize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	const statusBarHeight = 4
	const helpHeight = 2

	height := msg.Height - statusBarHeight - helpHeight

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
	var cmds []tea.Cmd

	switch keypress := msg.String(); keypress {
	case "J", "down":
		m.Messages, cmd = m.Messages.Update(message.NextEventMsg{})
		cmds = append(cmds, cmd)
		m.Timestamp, cmd = m.Timestamp.Update(timestamp.NextEventMsg{})
		cmds = append(cmds, cmd)
		return m, cmd
	case "K", "up":
		m.Messages, cmd = m.Messages.Update(message.PrevEventMsg{})
		cmds = append(cmds, cmd)
		m.Timestamp, cmd = m.Timestamp.Update(timestamp.PrevEventMsg{})
		cmds = append(cmds, cmd)
		return m, cmd
	case "enter":
		m.Messages, cmd = m.Messages.Update(
			message.ToggleCollapsedMsg{ToggleAll: false},
		)
		return m, cmd
	case "c":
		m.Messages, cmd = m.Messages.Update(
			message.ToggleCollapsedMsg{ToggleAll: true},
		)
		return m, cmd
	case "L":
		return m, m.loadMoreEvents()
	case "/":
		m.Timestamp, cmd = m.Timestamp.Update(msg)
		return m, cmd
	default:
		m.Timestamp, cmd = m.Timestamp.Update(msg)
		return m, cmd
	}
}

func (m *Model) loadMoreEvents() tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	ctx := context.Background()

	events := m.eventPaginator.NextPage(ctx)
	if events == nil {
		return nil
	}

	{ // update models with events
		m.Timestamp, cmd = m.Timestamp.Update(
			timestamp.LoadMoreEventsMsg(events),
		)
		cmds = append(cmds, cmd)

		m.Messages, cmd = m.Messages.Update(message.LoadMoreEventsMsg{
			AwsLogEvents: events,
			Collapsed:    true,
		})
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}
