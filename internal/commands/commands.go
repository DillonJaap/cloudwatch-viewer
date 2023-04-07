package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateViewPortContentMsg struct {
	Content []string
	YOffset int
}

type UpdateEventListItemsMsg struct {
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
