package logevent

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	numberOfEvents int
	selectedGroup  string
	selectedStream string
	selectedEvent  int
	help           help.Model
}

func New(
	timestampModel timestamp.Model,
	msg message.Model,
	initialGroup, initialStream string,
) Model {
	helpModel := help.New()
	helpModel.ShowAll = true

	model := Model{
		Timestamp:      timestampModel,
		Messages:       message.Model{},
		eventPaginator: nil,
		numberOfEvents: 0,
		selectedGroup:  initialGroup,
		selectedStream: "",
		selectedEvent:  0,
		help:           helpModel,
	}

	if initialGroup != "" {
		model, _ = model.Update(commands.UpdateStreamListItemsMsg{Group: initialGroup})
		model, _ = model.Update(commands.UpdateEventListItemsMsg{
			Group:  initialGroup,
			Stream: initialStream,
		})
	}

	return model
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
		m, cmd = m.updateEventItems()
		return m, cmd
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
	var cmds []tea.Cmd

	switch {
	case key.Matches(msg, keys.NextItem):
		if m.numberOfEvents-1 <= m.selectedEvent {
			return m, nil
		}
		m.selectedEvent += 1
		m.Messages, cmd = m.Messages.Update(message.NextEventMsg{
			Index: m.selectedEvent,
		})
		cmds = append(cmds, cmd)

		m.Timestamp, cmd = m.Timestamp.Update(timestamp.NextEventMsg{})
		cmds = append(cmds, cmd)

		return m, cmd
	case key.Matches(msg, keys.PrevItem):
		if m.selectedEvent <= 0 {
			return m, nil
		}
		m.selectedEvent -= 1
		m.Messages, cmd = m.Messages.Update(message.PrevEventMsg{
			Index: m.selectedEvent,
		})
		cmds = append(cmds, cmd)

		m.Timestamp, cmd = m.Timestamp.Update(timestamp.PrevEventMsg{})
		cmds = append(cmds, cmd)

		return m, cmd
	case key.Matches(msg, keys.Collapse):
		m.Messages, cmd = m.Messages.Update(
			message.ToggleCollapsedMsg{ToggleAll: false},
		)
		return m, cmd
	case key.Matches(msg, keys.CollapseAll):
		m.Messages, cmd = m.Messages.Update(
			message.ToggleCollapsedMsg{ToggleAll: true},
		)
		return m, cmd
	case key.Matches(msg, keys.Copy):
		m.Messages, cmd = m.Messages.Update(message.CopyMessage{})
		return m, cmd
	case key.Matches(msg, keys.LoadMore):
		return m, m.loadMoreEvents()
	case key.Matches(msg, keys.Reload):
		m, cmd = m.updateEventItems()
		return m, cmd
	case
		key.Matches(msg, keys.ScrollDown),
		key.Matches(msg, keys.ScrollUp),
		key.Matches(msg, keys.PageUp),
		key.Matches(msg, keys.PageDown),
		key.Matches(msg, keys.HalfPageUp),
		key.Matches(msg, keys.HalfPageDown):
		m.Messages, cmd = m.Messages.Update(msg)
		return m, cmd
	default:
		m.Timestamp, cmd = m.Timestamp.Update(msg)
		return m, cmd
	}
}

func (m Model) updateEventItems() (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// get a new paginator for our log group & stream
	paginator := event.New(
		context.Background(),
		m.selectedGroup,
		m.selectedStream,
		"",
	)
	m.eventPaginator = &paginator

	{ // reset data
		m.numberOfEvents = 0
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
	m.numberOfEvents += len(events)

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

func (m Model) HelpView() string {
	return m.help.View(keys)
}
