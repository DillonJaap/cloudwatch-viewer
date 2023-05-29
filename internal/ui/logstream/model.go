package logstream

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/cloudwatch/stream"
	"clviewer/internal/commands"
)

const listHeight = 14

var (
	titleStyle = lipgloss.
			NewStyle().
			Background(lipgloss.Color("98")).
			Foreground(lipgloss.Color("230")).
			PaddingLeft(1).
			PaddingRight(1)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type Model struct {
	List            list.Model
	SelectedStream  string
	currentGroup    string
	streamPaginator *stream.Paginator
}

func New(
	title string,
	initialGroup string,
) Model {
	streamList := list.New([]list.Item{}, &ItemDelegate{}, 0, 0)

	streamList.SetShowStatusBar(false)
	streamList.SetFilteringEnabled(true)
	streamList.SetShowHelp(false)

	streamList.Title = title
	streamList.Styles.Title = titleStyle
	streamList.Styles.PaginationStyle = paginationStyle

	model := Model{
		List:            streamList,
		SelectedStream:  "",
		currentGroup:    initialGroup,
		streamPaginator: &stream.Paginator{},
	}

	// initial group passed form cmd line arguments
	if initialGroup != "" {
		model.currentGroup = initialGroup
		model, _ = model.UpdateStreamItems()
	}
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO make other components also have independent height/width feilds
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "L":
			return m, m.loadMoreStreams()
		case "R":
			m, cmd := m.UpdateStreamItems()
			return m, cmd
		case "enter":
			if m.List.SettingFilter() {
				m.List, cmd = m.List.Update(msg)
				return m, cmd
			}

			i, ok := m.List.SelectedItem().(Item)
			if ok {
				m.SelectedStream = i.name
			}

			return m, commands.UpdateEventListItems(m.currentGroup, m.SelectedStream)
		}
	case commands.UpdateStreamListItemsMsg:
		m.currentGroup = msg.Group
		m, cmd = m.UpdateStreamItems()
		cmds = append(cmds, cmd)
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		PaddingRight(m.List.Width() - lipgloss.Width(m.List.View())).
		Render(m.List.View())
}

func (m Model) UpdateStreamItems() (Model, tea.Cmd) {
	ctx := context.Background()

	// reset list
	m.SelectedStream = ""
	m.List.ResetSelected()
	m.List.SetItems(nil)

	// get a new paginator for our log stream
	paginator := stream.New(ctx, m.currentGroup)
	m.streamPaginator = &paginator

	// itemList := GetLogStreamsAsItemList(groupPattern)
	// return m.List.SetItems(itemList)
	return m, m.loadMoreStreams()
}

func (m *Model) loadMoreStreams() tea.Cmd {
	ctx := context.Background()

	streams := m.streamPaginator.NextPage(ctx)
	if streams == nil {
		return nil
	}

	// Get streams into a formatted item list
	itemList := m.List.Items()
	itemList = append(itemList, GetLogStreamsAsItemList(streams)...)

	return m.List.SetItems(itemList)
}

func (m Model) HelpView() string {
	return m.List.Styles.HelpStyle.Render(m.List.Help.View(m.List))
}
