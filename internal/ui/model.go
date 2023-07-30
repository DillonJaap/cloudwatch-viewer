package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
	event "clviewer/internal/ui/logevent"
	"clviewer/internal/ui/logevent/message"
	"clviewer/internal/ui/logevent/timestamp"
	group "clviewer/internal/ui/loggroup"
	stream "clviewer/internal/ui/logstream"
	"clviewer/internal/ui/pages"
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
	groupListSelected = iota
	streamListSelected
	eventListSelected
	numWindows
)

const (
	groupPage = iota
	eventPage
)

type Model struct {
	paginator paginator.Model
	eventPage pages.Event
	groupPage pages.Group

	Width    int
	Height   int
	helpView string
	selected int
}

func New(ctx context.Context, initialGroup string) *Model {
	logGroup := group.New(
		"Log Groups",
		"/aws/lambda",
		initialGroup,
	)
	logStream := stream.New(
		"Log Streams",
		initialGroup,
	)

	initialLogstream := ""
	// initialize logstream only if logstream has been initalized with items
	if len(logStream.List.Items()) > 0 {
		initialLogstream = logStream.List.Items()[0].FilterValue()
	}

	logEvent := event.New(
		timestamp.New("Timestamps"),
		message.New("Log Messages", "..."),
		initialGroup,
		initialLogstream,
	)

	paginator := paginator.New()
	paginator.SetTotalPages(2)

	model := Model{
		eventPage: pages.Event{
			LogEvents:  logEvent,
			LogStreams: logStream,
			Focused:    0,
			Width:      0,
			Height:     0,
		},
		groupPage: pages.Group{
			Model: logGroup,
		},
		Width:     0,
		Height:    0,
		helpView:  "",
		selected:  eventListSelected,
		paginator: paginator,
	}
	return &model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	switch m.currentPage() {
	case groupPage:
		return m.groupPage.View()
	case eventPage:
		return m.eventPage.View()
	default:
		return ""
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h", "l":
			m.paginator, cmd = m.paginator.Update(msg)
			return m, cmd
		default:
			return m.updateCurrentPage(msg)
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m.updateWindowSizes()
	case commands.RedrawWindowsMsg:
		return m.updateWindowSizes()
	case commands.UpdateViewPortContentMsg:
		m.eventPage.Update(msg)
	}
	return m.updateCurrentPage(msg)
}

func (m *Model) updateWindowSizes() (*Model, tea.Cmd) {
	m, cmd := m.updateCurrentPage(
		tea.WindowSizeMsg{
			Width:  m.Width,
			Height: m.Height,
		},
	)
	return m, cmd
}

func (m *Model) updateCurrentPage(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m.currentPage() {
	case groupPage:
		m.groupPage, cmd = m.groupPage.Update(msg)
		return m, cmd
	case eventPage:
		m.eventPage, cmd = m.eventPage.Update(msg)
		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m Model) currentPage() int {
	return m.paginator.Page
}
