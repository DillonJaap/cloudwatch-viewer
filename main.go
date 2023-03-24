package main

import (
	"clviewer/internal/events"
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type eventListModel struct {
	events         list.Model
	cursorPosition int
	selected       []bool
}

func initialModel(ctx context.Context) eventListModel {
	cloudWatchEvents := events.GetEvents(ctx)

	selected := make([]bool, 5)
	for i := range selected {
		selected[i] = false
	}

	return eventListModel{
		events:   list.New(cloudWatchEvents, list.NewDefaultDelegate(), 0, 0),
		selected: selected,
	}
}

func (m eventListModel) Init() tea.Cmd {
	return nil
}

func (m eventListModel) View() string {
	items := m.events.Items()
	currentItem := items[m.events.Index()]

	eventList := docStyle.Render(m.events.View())
	message := docStyle.Render(currentItem.FilterValue())

	return lipgloss.JoinHorizontal(lipgloss.Center, eventList, message)
	//return docStyle.Render(m.events.View())
}

func (m eventListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.events.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.events, cmd = m.events.Update(msg)
	return m, cmd
}

func main() {
	ctx := context.TODO()
	p := tea.NewProgram(initialModel(ctx))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
