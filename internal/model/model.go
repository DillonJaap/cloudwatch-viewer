package model

import (
	"clviewer/internal/model/logevent"
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	eventsModel    logevent.Model
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

	return Model{
		eventsModel: logevent.DefaultModel(),
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
	//list := m.eventsModel.List
	//items := list.Items()
	//currentItem := items[list.Index()]

	//logGroupList := docStyle.Render(m.logGroupsModel.View())
	eventList := docStyle.Render(m.eventsModel.View())
	return eventList
	//message := docStyle.Render(currentItem.FilterValue())

	//return lipgloss.JoinHorizontal(lipgloss.Center, eventList)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m.eventsModel.Update(msg)
}
