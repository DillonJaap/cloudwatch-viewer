package pages

import (
	"log"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
	"clviewer/internal/ui/logevent"
	"clviewer/internal/ui/logstream"
)

const (
	logStreamsSelected = iota
	logEventsSelected
)

const (
	numWindows = 2
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

type Event struct {
	LogEvents  logevent.Model
	LogStreams logstream.Model
	Focused    int
	Width      int
	Height     int
}

func (m Event) Init() tea.Cmd {
	return nil
}

func (e Event) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			e = e.focusNext()
			return e, nil
		case "shift+tab":
			e = e.focusPrevious()
			return e, nil
		default:
			return e.updateKeyMsg(msg)
		}
	case tea.WindowSizeMsg:
		e.Width = msg.Width
		e.Height = msg.Height
		return e.updateWindowSizes()
	case commands.RedrawWindowsMsg:
		return e.updateWindowSizes()
	case commands.UpdateViewPortContentMsg:
		e.LogEvents.Update(msg)
	}

	return e.updateSubModels(msg)
}

func (e Event) View() string {
	logStreamList := e.LogStreams.View()
	logEventView := e.LogEvents.View()

	switch e.Focused {
	case logStreamsSelected:
		logStreamList = selectedModelStyle.Render(logStreamList)
		logEventView = modelStyle.Render(logEventView)
	case logEventsSelected:
		logEventView = selectedModelStyle.Render(logEventView)
		logStreamList = modelStyle.Render(logStreamList)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		logStreamList,
		logEventView,
	)
}

func (e Event) updateKeyMsg(msg tea.Msg) (Event, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch e.Focused {
	case logStreamsSelected:
		e.LogStreams, cmd = e.LogStreams.Update(msg)
	case logEventsSelected:
		e.LogEvents, cmd = e.LogEvents.Update(msg)
	}
	return e, cmd
}

func (e Event) updateWindowSizes() (Event, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	const borderMarginSize = 4 // subtract 4 for border
	const tuiBorder = 1

	height := e.Height - tuiBorder
	width := e.Width - tuiBorder

	logStreamListWidth := int(float64(width) / 3.0)
	e.LogStreams, cmd = e.LogStreams.Update(tea.WindowSizeMsg{
		Width:  logStreamListWidth,
		Height: height - borderMarginSize/2,
	})
	cmds = append(cmds, cmd)

	e.LogEvents, cmd = e.LogEvents.Update(tea.WindowSizeMsg{
		Width:  width - logStreamListWidth - borderMarginSize,
		Height: height - borderMarginSize,
	})
	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

func (e Event) updateSubModels(msg tea.Msg) (Event, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	e.LogStreams, cmd = e.LogStreams.Update(msg)
	cmds = append(cmds, cmd)

	e.LogEvents, cmd = e.LogEvents.Update(msg)
	cmds = append(cmds, cmd)

	log.Printf("%+v", e)
	return e, tea.Batch(cmds...)
}

func (e Event) focusNext() Event {
	e.Focused = (e.Focused + 1) % numWindows
	return e
}

func (e Event) focusPrevious() Event {
	e.Focused = int(math.Abs(float64((e.Focused - 1) % numWindows)))
	return e
}
