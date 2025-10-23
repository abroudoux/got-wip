package program

import "github.com/charmbracelet/huh"

func Confirm(msg string) bool {
	var confirmation bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title(msg).Value(&confirmation),
		),
	)

	err := form.Run()
	if err != nil {
		return false
	}

	return confirmation
}
