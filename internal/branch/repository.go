package branch

import (
	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func getRepository() (*repository, error) {
	gitRepository, err := git.PlainOpen(".")
	if err != nil {
		log.Warnf("Error opening git repository: %v", err)
		return nil, err
	}

	return gitRepository, nil
}

func getBranches(r *repository) ([]*branch, error) {
	branchIter, err := r.Branches()
	if err != nil {
		return nil, err
	}

	var branches []*branch

	err = branchIter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref)
		return nil
	})

	return branches, err
}

func getHead(r *repository) (*branch, error) {
	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	return head, nil
}

func findBranchByName(branchName string, r *repository) (*branch, error) {
	branches, err := getBranches(r)
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
