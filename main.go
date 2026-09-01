package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"go-choose-license/internal/license"
	"go-choose-license/internal/ui"
)

func main() {
	reg, err := license.LoadRegistry()
	if err != nil {
		fmt.Printf("Error initializing license registry: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.InitialModel(reg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running license chooser CLI: %v\n", err)
		os.Exit(1)
	}
}
