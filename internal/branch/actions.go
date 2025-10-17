package branch

import (
	"fmt"

	"github.com/abroudoux/got/internal/program"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func (r *repository) execAction(branchSelected *branch, action action) error {
	switch action {
	case actionNewBranch:
		return r.createNewBranch(branchSelected)
	case actionCopyName:
		return copyBranchName(branchSelected)
	case actionCheckout:
		return r.checkout(branchSelected)
	default:
		log.Info("Exiting..")
		return nil
	}
}

func copyBranchName(branch *branch) error {
	if clipboard.Unsupported {
		return fmt.Errorf("clipboard not supported on this plateform.")
	}

	err := clipboard.WriteAll(branch.Name().Short())
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Branch name %s copied to clipboard.", program.RenderElementSelected(branch.Name().Short())))
	return nil
}

func (r *repository) createNewBranch(branch *branch) error {
	if !r.isHead(branch) {
		log.Warn("You need to create a branch from the HEAD, move on it first.")
		return nil
	}

	for {
		// newBranchName, err := program.ReadInput("Enter the name of the new branch: ", "feat/amazing-feature")
		// if err != nil {
		// 	return err
		// }

		newBranchName := "hello"

		if r.isBranchNameAlreadyExists(newBranchName) {
			warnMsg := fmt.Sprintf("%s is already used, please choose another name.", program.RenderElementSelected(newBranchName))
			log.Warn(warnMsg)
			continue
		}

		newBranch := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+newBranchName), r.head.Hash())
		err := r.git.Storer.SetReference(newBranch)
		if err != nil {
			return nil
		}

		msgSuccessfullyCreated := fmt.Sprintf("New branch %s based on %s created.", program.RenderElementSelected(newBranchName), program.RenderElementSelected(r.head.Name().Short()))
		log.Info(msgSuccessfullyCreated)

		return nil
	}
}

func (r *repository) checkout(branch *branch) error {
	if r.isHead(branch) {
		log.Warn("You're already on the selected branch, please choose another one.")
		return nil
	}

	worktree, err := r.git.Worktree()
	if err != nil {
		return err
	}

	options := &git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branch.Name()),
	}

	err = worktree.Checkout(options)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Successfully checked out branch %s.", program.RenderElementSelected(branch.Name().Short()))
	log.Info(msg)

	return nil
}
