package branch

import (
	"fmt"

	"github.com/abroudoux/got/internal/program"
	tea "github.com/charmbracelet/bubbletea"
)

func selectAction(branchSelected *branch) (action, error) {
	p := tea.NewProgram(initialActionChoiceModel(branchSelected))
	m, err := p.Run()
	if err != nil {
		return actionExit, err
	}

	actionSelected := m.(actionChoice).actionSelected
	return actionSelected, nil
}

type action int

type actionChoice struct {
	actions        []action
	cursor         int
	actionSelected action
	branchSelected *branch
}

const (
	actionExit action = iota
	actionDelete
	actionMerge
	actionNewBranch
	actionCheckout
	actionRename
	actionPull
	actionCopyName
)

func (a action) String() string {
	return [...]string{
		"Exit",
		"Delete",
		"Merge",
		"New Branch",
		"Checkout",
		"Rename",
		"Pull",
		"Copy Name",
	}[a]
}

func initialActionChoiceModel(branch *branch) actionChoice {
	actions := []action{
		actionExit,
		actionDelete,
		actionMerge,
		actionNewBranch,
		actionCheckout,
		actionRename,
		actionPull,
		actionCopyName,
	}

	return actionChoice{
		actions:        actions,
		cursor:         len(actions) - 1,
		actionSelected: actionExit,
		branchSelected: branch,
	}
}

func (menu actionChoice) Init() tea.Cmd {
	return nil
}

func (menu actionChoice) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			menu.actionSelected = actionExit
			return menu, tea.Quit
		case "down":
			menu.cursor++
			if menu.cursor >= len(menu.actions) {
				menu.cursor = 0
			}
		case "up":
			menu.cursor--
			if menu.cursor < 0 {
				menu.cursor = len(menu.actions) - 1
			}
		case "enter":
			menu.actionSelected = menu.actions[menu.cursor]
			return menu, tea.Quit
		}
	}

	return menu, nil
}

func (menu actionChoice) View() string {
	s := "\033[H\033[2J\n"
	s += fmt.Sprintf("Choose an action for the branch %s:\n\n", program.RenderElementSelected(string(menu.branchSelected.Name().Short())))

	for i, action := range menu.actions {
		cursor := program.RenderCursor(menu.cursor == i)
		s += fmt.Sprintf("%s %s\n", cursor, program.RenderCurrentLine(action.String(), menu.cursor == i))
	}

	s += "\n"

	return s
}
