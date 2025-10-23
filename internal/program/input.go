package program

import "github.com/charmbracelet/huh"

func Input(msg string) (string, error) {
	var input string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title(msg).Value(&input),
		),
	)

	err := form.Run()
	if err != nil {
		return "", err
	}

	return input, nil
}
