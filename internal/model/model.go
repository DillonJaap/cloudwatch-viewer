package model

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	event "clviewer/internal/model/logeventlist"
	vp "clviewer/internal/model/logeventviewport"
	group "clviewer/internal/model/loggrouplist"
)

const listHeight = 15

var (
	docStyle          = lipgloss.NewStyle().Margin(1, 1)
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

var (
	eventListSelected = 1
	viewportSelected  = 2
)

type Model struct {
	eventsModel         event.Model
	viewportEventsModel vp.Model
	logGroupsModel      group.Model
	selected            int
}

func InitialModel(ctx context.Context) Model {
	items := []list.Item{
		group.Item("Ramen"),
		group.Item("Tomato Soup"),
		group.Item("Hamburgers"),
		group.Item("Cheeseburgers"),
		group.Item("Currywurst"),
		group.Item("Okonomiyaki"),
		group.Item("Pasta"),
		group.Item("Fillet Mignon"),
		group.Item("Caviar"),
		group.Item("Just Wine"),
	}
	const defaultWidth = 20

	logGroupList := list.New(items, group.LogGroupDelegate{}, defaultWidth, listHeight)
	logGroupList.SetShowStatusBar(false)
	logGroupList.SetFilteringEnabled(false)
	logGroupList.Title = "What do you want for dinner?"
	logGroupList.Styles.Title = titleStyle
	logGroupList.Styles.PaginationStyle = paginationStyle
	logGroupList.Styles.HelpStyle = helpStyle

	return Model{
		eventsModel:         event.DefaultModel(),
		viewportEventsModel: vp.InitialViewPortModel("test"),
		logGroupsModel:      group.Model{List: logGroupList, Choice: ""},
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
	case event.UpdateViewPort:
		eventValue := m.eventsModel.List.SelectedItem().FilterValue()
		m.viewportEventsModel.Events = eventValue
	}

	m.eventsModel, cmd = m.eventsModel.Update(msg)
	cmds = append(cmds, cmd)

	m.viewportEventsModel, cmd = m.viewportEventsModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
