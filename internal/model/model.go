package model

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/model/logevent"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var (
	eventListSelected = 1
	viewportSelected  = 2
)

type Model struct {
	eventsModel         logevent.Model
	viewportEventsModel viewPortModel
	logGroupsModel      LogGroupModel
	selected            int
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
		eventsModel:         logevent.DefaultModel(),
		viewportEventsModel: initialViewPortModel("test"),
		logGroupsModel:      LogGroupModel{list: logGroupList, choice: ""},
		selected:            eventListSelected,
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
	message := docStyle.Render(m.viewportEventsModel.View())

	return lipgloss.JoinHorizontal(lipgloss.Center, eventList, message)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.selected == eventListSelected {
			m.eventsModel, cmd = m.eventsModel.Update(msg)
			return m, cmd
		}
		if m.selected == viewportSelected {
			m.viewportEventsModel, cmd = m.viewportEventsModel.Update(msg)
			return m, cmd
		}
	default:
		eventValue := m.eventsModel.List.SelectedItem().FilterValue()
		m.viewportEventsModel.events = eventValue

		m.eventsModel, cmd = m.eventsModel.Update(msg)
		cmds = append(cmds, cmd)

		m.viewportEventsModel, cmd = m.viewportEventsModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
