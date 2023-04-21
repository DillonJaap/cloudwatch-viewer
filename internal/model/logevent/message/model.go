package message

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"clviewer/internal/commands"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.
			NewStyle().
			Background(lipgloss.Color("98")).
			Foreground(lipgloss.Color("230")).
			PaddingLeft(1).
			PaddingRight(1)
	}()

	infoStyle = func() lipgloss.Style {
		return titleStyle.Copy()
	}()

	lineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("98"))

	selectedItemStyleViewPort = lipgloss.NewStyle().Foreground(lipgloss.Color("127"))
)

type Model struct {
	Title    string
	Events   string
	Ready    bool
	Viewport viewport.Model
	messages []message
	index    int
}

type message struct {
	content    string
	collapsed  bool
	lineNumber int
}

func New(title string, events string) Model {
	return Model{
		Title:    title,
		Events:   events,
		Ready:    false,
		Viewport: viewport.Model{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type LoadMoreEventsMsg struct {
	AwsLogEvents []types.OutputLogEvent
	Collapsed    bool
}

type ResetMsg struct{}

type ToggleCollapsedMsg struct{}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.Viewport, cmd = m.Viewport.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.Ready {
			m.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.Viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.Viewport.SetContent(m.Events)
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.Viewport))
		}
	case commands.UpdateViewPortContentMsg:
		m.Viewport.SetContent(m.FormatList(msg.Content))
		// Center selected item in the viewport
		// TODO add updateViewPort and viewport scroll msg
		m.Viewport.SetYOffset(max(0, msg.YOffset-(m.Viewport.Height/2)))
		return m, nil
	case LoadMoreEventsMsg:
		m.messages = eventsToMessages(msg.AwsLogEvents, msg.Collapsed)
		m.Viewport.SetContent(m.renderContent())
		return m, nil
	case ToggleCollapsedMsg:
		collapseItems := true

		// if any are collapsed then don't set all to collapsed
		for k := range m.messages {
			if m.messages[k].collapsed {
				collapseItems = false
				break
			}
		}
		for k := range m.messages {
			m.messages[k].collapsed = collapseItems
		}

		m.Viewport.SetContent(m.renderContent())
		return m, nil
	}

	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) FormatList(list []string) string {
	var str string
	for _, msg := range list {
		str += msg + "\n"
	}
	return lipgloss.NewStyle().
		PaddingLeft(2).
		Render(str)
}

func (m Model) View() string {
	if !m.Ready {
		return "\n  Initializing..."
	}

	viewPortView := lipgloss.NewStyle().Margin(1, 1)

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.headerView(),
		viewPortView.Render(m.Viewport.View()),
		m.footerView(),
	)
}

func (m Model) headerView() string {
	title := titleStyle.Render(m.Title)
	line := lineStyle.Render(
		strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(title))),
	)
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	line := lineStyle.Render(strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m *Model) renderContent() string {
	var content string

	for i, event := range m.messages {
		formattedItem := FormatMessage(
			event.content,
			!event.collapsed,
		)

		if m.index == i {
			content += selectedItemStyleViewPort.Render(formattedItem) + "\n"
		} else {
			content += formattedItem + "\n"
		}

		// Set line number
		m.messages[i].lineNumber = lipgloss.Height(formattedItem)
	}
	return content
}

func FormatMessage(in string, formatAsJson bool) string {
	in = strings.ReplaceAll(in, "\t", " ")
	in = strings.ReplaceAll(in, "\n", " ")

	if in[0] == '{' && formatAsJson {
		return formatJson(in)
	}
	return in
}

func formatJson(in string) string {
	var obj map[string]interface{}
	json.Unmarshal([]byte(in), &obj)

	f := colorjson.NewFormatter()
	f.Indent = 2

	s, _ := f.Marshal(obj)
	return string(s)
}

func eventsToMessages(logEvents []types.OutputLogEvent, collaped bool) []message {
	var events []message
	for k := range logEvents {
		events = append(
			events,
			message{
				content:    aws.ToString(logEvents[k].Message),
				collapsed:  collaped,
				lineNumber: k, // TODO should line number start at one?
			},
		)
	}

	return events
}
