package timestamp

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
) Model {
	eventList := list.New([]list.Item{}, &ItemDelegate{}, 0, 0)

	eventList.SetShowStatusBar(false)
	eventList.SetFilteringEnabled(true)
	eventList.SetShowHelp(false)

	eventList.Title = title
	eventList.Styles.Title = titleStyle
	eventList.Styles.PaginationStyle = paginationStyle
	eventList.Styles.HelpStyle = helpStyle

	return Model{
		List:         eventList,
		Choice:       "",
		ItemMetaData: []ItemMetaData{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height)
		return m, nil
	case tea.KeyMsg:
		return m.updateKeyMsg(msg)
	case commands.UpdateEventListItemsMsg:
		m, cmd = m.updateEventListItems(msg.Group, msg.Stream, false)
		cmds = append(cmds, cmd)
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.List.View()
}

// updateEventListItems updates the event lists with new messages from the new
// group/stream and refreshes the viewport
func (m Model) updateEventListItems(
	groupPattern string,
	streamPrefix string,
	collapsed bool,
) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	// reset list
	m.Choice = ""

	// update item list
	itemList := GetLogEventsAsItemList(groupPattern, streamPrefix)
	itemList = formatList(itemList, false)
	cmd = m.List.SetItems(itemList)
	cmds = append(cmds, cmd)

	// update item meta data
	metaData := make([]ItemMetaData, len(itemList))
	for i := range metaData {
		metaData[i].Collapsed = collapsed
	}
	m.ItemMetaData = metaData

	// update viewport
	cmd = commands.UpdateViewPort(
		m.getItemListAsStringArray(),
		m.getLineNumber(),
	)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// updateKeyMsg updates model based on the tea.KeyMsg
func (m Model) updateKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	index := m.List.Index()
	switch keypress := msg.String(); keypress {
	// Toggle Collapse on Item
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
	// Toggle Collapse all
	case "c":
		m.toggleCollapseAll()
		cmd = commands.UpdateViewPort(
			m.getItemListAsStringArray(),
			m.ItemMetaData[index].lineNum,
		)
		return m, cmd
	case "J":
		m.List.CursorDown()
		cmd = commands.UpdateViewPort(
			m.getItemListAsStringArray(),
			m.getLineNumber(),
		)
		return m, cmd
	case "K":
		m.List.CursorUp()
		cmd = commands.UpdateViewPort(
			m.getItemListAsStringArray(),
			m.getLineNumber(),
		)
		return m, cmd
	// all other keystrokes get handled by the list Model
	// and then the viewport gets updated
	default:
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)

		cmd = commands.UpdateViewPort(
			m.getItemListAsStringArray(),
			m.getLineNumber(),
		)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

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

func (m Model) getLineNumber() int {
	lineNum := 0
	if len(m.ItemMetaData) > m.List.Index()+1 {
		lineNum = m.ItemMetaData[m.List.Index()].lineNum
	}
	return lineNum
}
