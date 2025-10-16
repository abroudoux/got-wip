package branch

import (
	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
)

func getCurrentGitRepository() (*Repository, error) {
	currentDir, err := git.PlainOpen(".")
	if err != nil {
		log.Warnf("Error opening git repository: %v", err)
		return nil, err
	}

	var repository = NewRepository(currentDir)
	repository.head, err = repository.getHead()
	if err != nil {
		log.Warnf("Error getting HEAD: %v", err)
		return nil, err
	}

	repository.branches, err = repository.getBranches()
	if err != nil {
		log.Warnf("Error getting branches: %v", err)
		return nil, err
	}

	return repository, nil
}

func (repository *Repository) getHead() (*Branch, error) {
	head, err := repository.git.Head()
	if err != nil {
		return nil, err
	}

	return head, nil
}

func (repository *Repository) getBranches() ([]*Branch, error) {
	branchIter, err := repository.git.Branches()
	if err != nil {
		return nil, err
	}

	var branches []*Branch
	err = branchIter.ForEach(func(ref *Branch) error {
		branches = append(branches, ref)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return branches, nil
}