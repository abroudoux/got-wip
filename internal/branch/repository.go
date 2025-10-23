package branch

import (
	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type branch = plumbing.Reference
type gitRepository = git.Repository

type repository struct {
	gitRepository
}

func newRepository(gitRepository *gitRepository) *repository {
	return &repository{
		*gitRepository,
	}
}

func getRepository() (*repository, error) {
	gitRepository, err := git.PlainOpen(".")
	if err != nil {
		log.Warnf("Error opening git repository: %v", err)
		return nil, err
	}

	return newRepository(gitRepository), nil
}

func (r *repository) getBranches() ([]*branch, error) {
	branchIter, err := r.Branches()
	if err != nil {
		return nil, err
	}

	var branches []*branch

	err = branchIter.ForEach(func(ref *branch) error {
		branches = append(branches, ref)
		return nil
	})

	return branches, err
}

func (r *repository) findBranchByName(branchName string) (*branch, error) {
	branches, err := r.getBranches()
	if err != nil {
		return nil, err
	}

	for _, b := range branches {
		if b.Name().Short() == branchName {
			return b, nil
		}
	}

	return nil, nil
}

func (r *repository) isHead(branch *branch) bool {
	head, err := r.Head()
	if err != nil {
		return false
	}

	return head.Name().Short() == branch.Name().Short()
}

func (r *repository) isBranchNameAlreadyExists(branchName string) bool {
	branches, err := r.getBranches()
	if err != nil {
		return false
	}

	for _, b := range branches {
		if b.Name().Short() == branchName {
			return true
		}
	}

	return false
}
