package util

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func As[T any]() func(any, tea.Cmd) (T, tea.Cmd) {
	return func(in any, cmd tea.Cmd) (T, tea.Cmd) {
		model, ok := in.(T)
		if !ok {
			log.Fatalf("expected model to be of a different type")
		}
		return model, cmd
	}
}
