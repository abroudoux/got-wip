package branch

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type branch = plumbing.Reference
type gitRepository = git.Repository
type remote = git.Remote

type repository struct {
	git      *gitRepository
	branches []*branch
	head     *branch
	remotes  []*remote
}

type branchChoice struct {
	head           *branch
	branches       []*branch
	cursor         int
	branchSelected *branch
}

type action int

type actionChoice struct {
	actions        []action
	cursor         int
	actionSelected action
	branchSelected *branch
}
