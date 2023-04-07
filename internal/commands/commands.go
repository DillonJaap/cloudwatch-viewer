package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateViewPortContentMsg struct {
	Content []string
	YOffset int
}

type UpdateEventListItemsMsg struct {
	Group  string
	Stream string
}

type UpdateStreamListItemsMsg struct {
	Group string
}

func UpdateViewPort(content []string, yOffset int) tea.Cmd {
	return func() tea.Msg {
		return UpdateViewPortContentMsg{
			Content: content,
			YOffset: yOffset,
		}
	}
}

func UpdateEventListItems(group string) tea.Cmd {
	return func() tea.Msg {
		return UpdateEventListItemsMsg{
			Group: group,
		}
	}
}

func UpdateStreamListItems(group string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStreamListItemsMsg{
			Group: group,
		}
	}
}
