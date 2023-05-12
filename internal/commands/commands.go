package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateViewPortContentMsg struct {
	Content []string
	YOffset int
}

func UpdateViewPort(content []string, yOffset int) tea.Cmd {
	return func() tea.Msg {
		return UpdateViewPortContentMsg{
			Content: content,
			YOffset: yOffset,
		}
	}
}

type UpdateEventListItemsMsg struct {
	Group  string
	Stream string
}

func UpdateEventListItems(group string, stream string) tea.Cmd {
	return func() tea.Msg {
		return UpdateEventListItemsMsg{
			Group:  group,
			Stream: stream,
		}
	}
}

type UpdateStreamListItemsMsg struct {
	Group string
}

func UpdateStreamListItems(group string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStreamListItemsMsg{
			Group: group,
		}
	}
}

type RedrawWindowsMsg struct{}

func RedrawWindows() tea.Cmd {
	return func() tea.Msg {
		return RedrawWindowsMsg{}
	}
}
