package logeventviewport

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	/*s.Title = lipgloss.NewStyle().
	Background(lipgloss.Color("62")).
	Foreground(lipgloss.Color("230")).
	Padding(0, 1)
	*/

	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.
			NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			PaddingLeft(1).
			PaddingRight(1).
			Margin(1)
	}()

	infoStyle = func() lipgloss.Style {
		return titleStyle.Copy()
	}()

	lineStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
	viewPortStyle = lipgloss.NewStyle().Padding(2)
)

type Model struct {
	Events   string
	Ready    bool
	Viewport viewport.Model
}

func New(events string) Model {
	return Model{
		Events: events,
		Ready:  false,
		Viewport: viewport.Model{
			Style: viewPortStyle,
		},
	}
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
	case tea.KeyMsg:
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
	}

	// Handle keyboard and mouse events in the viewport
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
	title := titleStyle.Render("Log Message")
	line := lineStyle.Render(
		strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(title))),
	)
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	line := lineStyle.Render(
		strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(info))),
	)
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
