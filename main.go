package main

import (
	"clviewer/internal/model"
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx := context.TODO()
	p := tea.NewProgram(model.InitialModel(ctx), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
