package branch

import "fmt"

func SelectBranch() error {
	repository, err := getCurrentGitRepository();
	if err != nil {
		return err
	}

	for _, b := range repository.branches {
		fmt.Println(b.Name())
	}

	return nil
}