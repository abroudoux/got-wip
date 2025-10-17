package branch

import "github.com/charmbracelet/log"

func Run() error {
	repository, err := getCurrentGitRepository()
	if err != nil {
		log.Errorf("Error while getting current git repository: %v", err)
		return err
	}

	branch, err := repository.selectBranch()
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

	err = repository.execAction(branch, action)
	if err != nil {
		log.Errorf("Error while executing action: %v", err)
		return err
	}

	return nil
}
