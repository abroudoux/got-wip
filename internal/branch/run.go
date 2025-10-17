package branch

import "github.com/charmbracelet/log"

func Run() error {
	repository, err := getRepository()
	if err != nil {
		log.Errorf("Error while getting current git repository: %v", err)
		return err
	}

	branch, err := selectBranch(repository)
	if branch == nil {
		return nil
	}
	if err != nil {
		log.Errorf("Error while selecting branch: %v", err)
		return err
	}

	action, err := selectAction(branch)
	if err != nil {
		log.Errorf("Error while selecting action: %v", err)
		return err
	}

	err = execAction(branch, action, repository)
	if err != nil {
		log.Errorf("Error while executing action: %v", err)
		return err
	}

	return nil
}
