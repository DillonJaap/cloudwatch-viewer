package pages

import (
	group "clviewer/internal/ui/loggroup"

	tea "github.com/charmbracelet/bubbletea"
)

type Group struct {
	group.Model
}

func (g Group) Init() tea.Cmd {
	return g.Model.Init()
}

func (g Group) Update(msg tea.Msg) (Group, tea.Cmd) {
	var cmd tea.Cmd
	g.Model, cmd = g.Model.Update(msg)
	return g, cmd
}

func (g Group) View() string {
	return g.Model.View()
}
