package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"job-applier-tui/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
