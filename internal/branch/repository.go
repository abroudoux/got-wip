package branch

import (
	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
)

func newRepository(gitRepository *gitRepository) *repository {
	return &repository{
		git: gitRepository,
	}
}

func getCurrentGitRepository() (*repository, error) {
	currentDir, err := git.PlainOpen(".")
	if err != nil {
		log.Warnf("Error opening git repository: %v", err)
		return nil, err
	}

	var repository = newRepository(currentDir)
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

func (repository *repository) getHead() (*branch, error) {
	head, err := repository.git.Head()
	if err != nil {
		return nil, err
	}

	return head, nil
}

func (repository *repository) getBranches() ([]*branch, error) {
	branchIter, err := repository.git.Branches()
	if err != nil {
		return nil, err
	}

	var branches []*branch
	err = branchIter.ForEach(func(ref *branch) error {
		branches = append(branches, ref)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *repository) findBranchByName(name string) *branch {
	r.branches, _ = r.getBranches()

	for _, branch := range r.branches {
		if branch.Name().Short() == name {
			return branch
		}
	}

	return nil
}
