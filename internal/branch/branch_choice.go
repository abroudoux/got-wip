package branch

import (
	"fmt"

	"github.com/abroudoux/got/internal/program"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func selectBranch(r *repository) (*branch, error) {
	branches, err := getBranches(r)
	if err != nil {
		return nil, err
	}

	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	p := tea.NewProgram(initialBranchChoiceModel(branches, head))
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	branchSelected := m.(branchChoice).branchSelected
	return branchSelected, nil
}

func initialBranchChoiceModel(branches []*branch, head *branch) branchChoice {
	return branchChoice{
		head:           head,
		branches:       branches,
		cursor:         len(branches) - 1,
		branchSelected: nil,
	}
}

func (menu branchChoice) Init() tea.Cmd {
	return nil
}

func (menu branchChoice) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			menu.branchSelected = nil
			log.Info("Program exited.")
			return menu, tea.Quit
		case "down":
			menu.cursor++
			if menu.cursor >= len(menu.branches) {
				menu.cursor = 0
			}
		case "up":
			menu.cursor--
			if menu.cursor < 0 {
				menu.cursor = len(menu.branches) - 1
			}
		case "enter":
			menu.branchSelected = menu.branches[menu.cursor]
			return menu, tea.Quit
		}
	}

	return menu, nil
}

func (menu branchChoice) View() string {
	s := "\033[H\033[2J\n"
	s += "Choose a branch:\n\n"

	for i, branch := range menu.branches {
		cursor := program.RenderCursor(menu.cursor == i)

		if branch.Name().Short() == menu.head.Name().Short() {
			branchName := "* " + branch.Name().Short()
			s += fmt.Sprintf("%s %s\n", cursor, program.RenderCurrentLine(branchName, menu.cursor == i))
		} else {
			branchName := "  " + branch.Name().Short()
			s += fmt.Sprintf("%s %s\n", cursor, program.RenderCurrentLine(branchName, menu.cursor == i))
		}
	}

	s += "\n"

	return s
}
