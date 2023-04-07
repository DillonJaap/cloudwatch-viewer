package logeventlist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
)

var (
	titleStyle = lipgloss.
			NewStyle().
			Background(lipgloss.Color("98")).
			Foreground(lipgloss.Color("230")).
			PaddingLeft(1).
			PaddingRight(1)

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

func New(
	title string,
	collapsed bool,
) *Model {
	eventList := list.New([]list.Item{}, &ItemDelegate{}, 0, 0)

	eventList.SetShowStatusBar(false)
	eventList.SetFilteringEnabled(true)
	eventList.SetShowHelp(false)

	eventList.Title = title
	eventList.Styles.Title = titleStyle
	eventList.Styles.PaginationStyle = paginationStyle
	eventList.Styles.HelpStyle = helpStyle

	return &Model{
		List:         eventList,
		Choice:       "",
		ItemMetaData: []ItemMetaData{},
	}
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
				m.getItemListAsStringArray(),
				m.ItemMetaData[index].lineNum,
			)
			return m, cmd
		case "c":
			m.toggleCollapseAll()
			cmd = commands.UpdateViewPort(
				m.getItemListAsStringArray(),
				m.ItemMetaData[index].lineNum,
			)
			return m, cmd
		default:
			m.List, cmd = m.List.Update(msg)

			cmd = commands.UpdateViewPort(
				m.getItemListAsStringArray(),
				m.getLineNumber(),
			)
			cmds = append(cmds, cmd)

			return m, tea.Batch(cmds...)
		}
	case commands.UpdateEventListItemsMsg:
		cmd = m.UpdateEventItems(msg.Group, msg.Stream, false)
		cmds = append(cmds, cmd)
		cmd = commands.UpdateViewPort(
			m.getItemListAsStringArray(),
			m.getLineNumber(),
		)
		cmds = append(cmds, cmd)
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.List.View()
}

func (m *Model) UpdateEventItems(
	groupPattern string,
	streamPrefix string,
	collapsed bool,
) tea.Cmd {
	itemList := GetLogEventsAsItemList(groupPattern, streamPrefix)
	itemList = formatList(itemList, false)
	cmd := m.List.SetItems(itemList)

	metaData := make([]ItemMetaData, len(itemList))
	for i := range metaData {
		metaData[i].Collapsed = collapsed
	}
	m.ItemMetaData = metaData

	return cmd
}

func (m *Model) getItemListAsStringArray() []string {
	var list []string
	height := 0

	for index, item := range m.List.Items() {
		formattedItem := FormatMessage(
			item.FilterValue(),
			!m.ItemMetaData[index].Collapsed,
		)

		if m.List.Index() == index {
			list = append(list, selectedItemStyleViewPort.Render(formattedItem))
		} else {
			list = append(list, formattedItem)
		}

		height += lipgloss.Height(formattedItem)
		m.ItemMetaData[index].lineNum = height
	}
	return list
}

func (m *Model) toggleCollapseAll() {
	collapseItems := true

	for k := range m.ItemMetaData {
		if m.ItemMetaData[k].Collapsed {
			collapseItems = false
			break
		}
	}

	for k := range m.ItemMetaData {
		m.ItemMetaData[k].Collapsed = collapseItems
	}
}

func (m *Model) getLineNumber() int {
	lineNum := 0
	if len(m.ItemMetaData) > 0 {
		lineNum = m.ItemMetaData[m.List.Index()].lineNum
	}
	return lineNum
}
