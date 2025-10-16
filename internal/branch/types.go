package branch

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Branch = plumbing.Reference
type GitRepository = git.Repository

type Repository struct {
	git *GitRepository
	branches []*Branch
	head    *Branch
}

func NewRepository(gitRepository *GitRepository) *Repository {
	return &Repository{
		git: gitRepository,
	}
}

type BranchChoice struct {
	head *Branch
	branches []*Branch
	cursor int
	branchSelected *Branch
}