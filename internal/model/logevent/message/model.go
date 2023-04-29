package message

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
)

type Model struct {
	Title         string
	Events        string
	Ready         bool
	Viewport      viewport.Model
	messages      []message
	selectedEvent int
}

type message struct {
	content    string
	collapsed  bool
	lineNumber int
}

func New(title string, events string) Model {
	return Model{
		Title:         title,
		Events:        events,
		Ready:         false,
		Viewport:      viewport.Model{},
		messages:      []message{},
		selectedEvent: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type LoadMoreEventsMsg struct {
	AwsLogEvents []types.OutputLogEvent
	Collapsed    bool
}

type ToggleCollapsedMsg struct {
	ToggleAll bool
}

type ResetMsg struct{}

type NextEventMsg struct{ Index int }

type PrevEventMsg struct{ Index int }

type CopyMessage struct{}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "J" || msg.String() == "shift+down" {
			m.Viewport.LineDown(3)
		}
		if msg.String() == "K" || msg.String() == "shift+up" {
			m.Viewport.LineUp(3)
		}
		m.Viewport, cmd = m.Viewport.Update(msg)
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.Ready {
			m.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.Viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.Viewport))
		}
		return m, tea.Batch(cmds...)
	case ResetMsg:
		m.selectedEvent = 0
		m.messages = []message{}
	case NextEventMsg:
		m.selectedEvent = msg.Index
		m.centerViewOnItem()
	case PrevEventMsg:
		m.selectedEvent = msg.Index
		m.centerViewOnItem()
	case LoadMoreEventsMsg:
		m.messages = append(
			m.messages,
			eventsToMessages(msg.AwsLogEvents, msg.Collapsed)...,
		)
	case CopyMessage:
		eventMsg := m.messages[m.selectedEvent].content
		collapsed := m.messages[m.selectedEvent].collapsed

		eventMsg = FormatMessage(eventMsg, !collapsed)

		if err := clipboard.WriteAll(eventMsg); err != nil {
			log.Printf("error with clipboard: %s", err)
		}
		return m, nil
	case ToggleCollapsedMsg:
		// break if no messages have been set
		if len(m.messages) == 0 {
			break
		}

		// Toggle one item
		if !msg.ToggleAll {
			m.messages[m.selectedEvent].collapsed = !m.messages[m.selectedEvent].collapsed
			break
		}

		{ // Toggle all Items
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

		}
	}

	m.Viewport.SetContent(m.renderContent())
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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

func (m *Model) centerViewOnItem() {
	itemHeightOffset := lipgloss.Height(FormatMessage(
		m.messages[m.selectedEvent].content,
		!m.messages[m.selectedEvent].collapsed,
	))

	if itemHeightOffset > m.Viewport.Height {
		m.Viewport.SetYOffset(max(
			0,
			m.messages[m.selectedEvent].lineNumber-2, // TODO use const?
		))
		return
	}

	m.Viewport.SetYOffset(max(
		0,
		m.messages[m.selectedEvent].lineNumber-(m.Viewport.Height/2)+(itemHeightOffset/2),
	))
}

func (m *Model) renderContent() string {
	var content string

	for i, event := range m.messages {
		formattedItem := FormatMessage(
			event.content,
			!event.collapsed,
		)

		// Style if item is unselected
		style := lipgloss.NewStyle().
			BorderLeft(false).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeftForeground(lipgloss.Color("237")).
			PaddingLeft(3).
			PaddingRight(m.Viewport.Width - lipgloss.Width(formattedItem) - 3)

		// Style if item is selected
		if m.selectedEvent == i {
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("127")).
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeftForeground(lipgloss.Color("127")).
				PaddingLeft(3).
				PaddingRight(m.Viewport.Width - lipgloss.Width(formattedItem) - 4)
		}

		// format collapsed items
		if lipgloss.Height(formattedItem) == 1 {
			bgColor := "232"
			if i%2 == 0 {
				bgColor = "235"
			}
			formattedItem = style.
				Background(lipgloss.Color(bgColor)).
				Bold(true).
				Render(formattedItem)
		} else {
			// format expanded items
			formattedItem = style.Render(formattedItem)
		}

		// Set line number
		m.messages[i].lineNumber = lipgloss.Height(content) + 1

		content += formattedItem + "\n"

	}
	return content
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
