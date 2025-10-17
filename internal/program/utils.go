package program

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func clearScreen() tea.Cmd {
	return func() tea.Msg {
		fmt.Print("\033[H\033[2J")
		return nil
	}
}
