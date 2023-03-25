package model

import (
	"clviewer/internal/events"
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	eventsModel    list.Model
	logGroupsModel LogGroupModel
}

func InitialModel(ctx context.Context) Model {
	items := []list.Item{
		item("Ramen"),
		item("Tomato Soup"),
		item("Hamburgers"),
		item("Cheeseburgers"),
		item("Currywurst"),
		item("Okonomiyaki"),
		item("Pasta"),
		item("Fillet Mignon"),
		item("Caviar"),
		item("Just Wine"),
	}
	const defaultWidth = 20

	logGroupList := list.New(items, logGroupDelegate{}, defaultWidth, listHeight)
	logGroupList.SetShowStatusBar(false)
	logGroupList.SetFilteringEnabled(false)
	logGroupList.Title = "What do you want for dinner?"
	logGroupList.Styles.Title = titleStyle
	logGroupList.Styles.PaginationStyle = paginationStyle
	logGroupList.Styles.HelpStyle = helpStyle

	cloudWatchEvents := events.GetEvents(ctx)

	return Model{
		eventsModel: list.New(cloudWatchEvents, list.NewDefaultDelegate(), 0, 0),
		logGroupsModel: LogGroupModel{
			list:   logGroupList,
			choice: "",
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	items := m.eventsModel.Items()
	currentItem := items[m.eventsModel.Index()]

	//logGroupList := docStyle.Render(m.logGroupsModel.View())
	eventList := docStyle.Render(m.eventsModel.View())
	message := docStyle.Render(currentItem.FilterValue())

	return lipgloss.JoinHorizontal(lipgloss.Center, eventList, message)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.eventsModel.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.eventsModel, cmd = m.eventsModel.Update(msg)
	return m, cmd
}
