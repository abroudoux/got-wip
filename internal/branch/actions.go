package branch

import (
	"errors"
	"fmt"

	"github.com/abroudoux/got/internal/program"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func (r *repository) execAction(branchSelected *branch, action action) error {
	switch action {
	case actionDelete:
		// return r.delete(branchSelected)
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
		Keep:   true,
	}

	err = worktree.Checkout(options)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Successfully checked out branch %s.", program.RenderElementSelected(branch.Name().Short()))
	log.Info(msg)

	return nil
}

func (r *repository) createNewBranch(branch *branch) error {
	if !r.isHead(branch) {
		log.Warn("You need to create a branch from the HEAD, move on it first.")
		return nil
	}

	var newBranchName string
	var checkoutOnNewBranch bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter the name of the new branch:").Value(&newBranchName).Validate(
				func(str string) error {
					if r.isBranchNameAlreadyExists(str) {
						return errors.New("Branch name already exists, please choose another name.")
					}
					return nil
				}).Placeholder("feature/my-new-branch"),
			huh.NewConfirm().Title("Do you want to checkout on?").Value(&checkoutOnNewBranch),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	newBranch := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+newBranchName), r.head.Hash())
	err = r.git.Storer.SetReference(newBranch)
	if err != nil {
		return nil
	}

	msgSuccessfullyCreated := fmt.Sprintf("New branch %s based on %s created.", program.RenderElementSelected(newBranchName), program.RenderElementSelected(r.head.Name().Short()))
	log.Info(msgSuccessfullyCreated)

	if checkoutOnNewBranch {
		branchCreated := r.findBranchByName(newBranch.Name().Short())
		if branchCreated == nil {
			return fmt.Errorf("branch %s not found after creation", newBranch.Name().Short())
		}

		err := r.checkout(branchCreated)
		if err != nil {
			return err
		}
	}

	return nil
}
