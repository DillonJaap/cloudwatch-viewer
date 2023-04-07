package logeventlist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	selectedItemStyleViewPort = lipgloss.NewStyle().Foreground(lipgloss.Color("127"))
)

var _ tea.Model = &Model{}

type Model struct {
	List         list.Model
	Choice       string
	ItemMetaData []ItemMetaData
}

type ItemMetaData struct {
	Collapsed bool
	lineNum   int
}

func New(groupPattern string, collapsed bool) *Model {
	itemList := GetLogEventsAsItemList(groupPattern)
	itemList = formatList(itemList, false)

	eventList := list.New(itemList, &ItemDelegate{}, 0, 0)
	eventList.SetShowStatusBar(false)
	eventList.SetFilteringEnabled(true)
	eventList.Title = "Timestamp"
	eventList.Styles.PaginationStyle = paginationStyle
	eventList.Styles.HelpStyle = helpStyle
	eventList.Styles.Title = lipgloss.
		NewStyle().
		Background(lipgloss.Color("98")).
		Foreground(lipgloss.Color("230")).
		PaddingLeft(1).
		PaddingRight(1)

	metaData := make([]ItemMetaData, len(itemList))
	for i := range metaData {
		metaData[i].Collapsed = collapsed
	}

	return &Model{
		List:         eventList,
		Choice:       "",
		ItemMetaData: metaData,
	}
}

func (m *Model) GetEventItems(groupPattern string, collapsed bool) tea.Cmd {
	itemList := GetLogEventsAsItemList(groupPattern)
	itemList = formatList(itemList, false)
	cmd := m.List.SetItems(itemList)

	metaData := make([]ItemMetaData, len(itemList))
	for i := range metaData {
		metaData[i].Collapsed = collapsed
	}
	m.ItemMetaData = metaData

	return cmd
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	index := m.List.Index()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			item, ok := m.List.SelectedItem().(Item)
			if ok {
				m.Choice = item.Message
			}
			m.ItemMetaData[index].Collapsed = !m.ItemMetaData[index].Collapsed
			cmd = commands.UpdateViewPort(
				m.getItemListAsString(),
				m.ItemMetaData[index].lineNum,
			)
			return m, cmd
		default:
			m.List, cmd = m.List.Update(msg)
			cmd = commands.UpdateViewPort(
				m.getItemListAsString(),
				m.ItemMetaData[index].lineNum,
			)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	case commands.UpdateEventListItemsMsg:
		m.GetEventItems(msg.Group, false)
		cmds = append(cmds, cmd)
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.List.View()
}

func (m *Model) getItemListAsString() string {
	var list string

	for index, item := range m.List.Items() {
		formattedItem := FormatMessage(
			item.FilterValue(),
			!m.ItemMetaData[index].Collapsed,
		)

		if m.List.Index() == index {
			list += selectedItemStyleViewPort.Render(formattedItem) + "\n"
		} else {
			list += formattedItem + "\n"
		}

		// TODO could optimize this by keeping track of the current list height,
		// instead of recalculating it everytime
		m.ItemMetaData[index].lineNum = lipgloss.Height(list)
	}
	return list
}
