package branch

func Run() error {
	repository, err := getCurrentGitRepository()
	if err != nil {
		return err
	}

	branch, err := repository.selectBranch()
	if err != nil {
		return err
	}

	action, err := selectAction(branch)
	if err != nil {
		return err
	}

	err = repository.execAction(branch, action)
	if err != nil {
		return err
	}

	return nil
}
