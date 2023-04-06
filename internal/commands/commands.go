package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateViewPortContentMsg struct {
	Content string
}

type UpdateEventListItemsMsg struct {
	Group string
}

func UpdateViewPort(content string) tea.Cmd {
	return func() tea.Msg {
		return UpdateViewPortContentMsg{
			Content: content,
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
