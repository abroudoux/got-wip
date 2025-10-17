package program

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func ReadInput(msg string, placeholder string) (string, error) {
	p := tea.NewProgram(initialModel(msg, placeholder), tea.WithOutput(os.Stdout))
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	finalModel := m.(model)
	return finalModel.inputValue, nil
}

type model struct {
	prompt     string
	textInput  textinput.Model
	inputValue string
	cancelled  bool
	done       bool
}

func initialModel(msg string, placeholder string) model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return model{
		prompt:    msg,
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, clearScreen())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.done = true
			m.inputValue = m.textInput.Value()
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			m.done = true
			m.cancelled = true
			log.Info("Program exited..")
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.done {
		fmt.Print("\033[H\033[2J")
		return ""
	}

	return fmt.Sprintf(
		"\033[H\033[2J%s\n\n%s\n\n(Enter to confirm, Esc to cancel)",
		m.prompt,
		m.textInput.View(),
	)
}
