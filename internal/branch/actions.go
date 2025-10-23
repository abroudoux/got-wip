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

func execAction(branchSelected *branch, action action, r *repository) error {
	switch action {
	case actionDelete:
		return r.delete(branchSelected)
	case actionNewBranch:
		return r.createNewBranch(branchSelected)
	case actionCopyName:
		return copyBranchName(branchSelected)
	case actionCheckout:
		return r.checkout(branchSelected)
	case actionMerge:
		return r.merge(branchSelected)
	case actionPull:
		return r.pull(branchSelected)
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

	worktree, err := r.Worktree()
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
				func(input string) error {
					if r.isBranchNameAlreadyExists(input) {
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

	head, err := r.Head()
	if err != nil {
		return err
	}

	newBranch := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/"+newBranchName), head.Hash())
	err = r.Storer.SetReference(newBranch)
	if err != nil {
		return nil
	}

	msgSuccessfullyCreated := fmt.Sprintf("New branch %s based on %s created.", program.RenderElementSelected(newBranchName), program.RenderElementSelected(head.Name().Short()))
	log.Info(msgSuccessfullyCreated)

	if !checkoutOnNewBranch {
		return nil
	}

	branchCreated, err := r.findBranchByName(newBranch.Name().Short())
	if err != nil {
		return err
	}
	if branchCreated == nil {
		return fmt.Errorf("branch %s not found after creation", newBranch.Name().Short())
	}

	err = r.checkout(branchCreated)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) delete(branch *branch) error {
	if r.isHead(branch) {
		log.Warn("You cannot delete the branch you're currently on, please checkout on another branch first.")
		return nil
	}

	confirmDelete := program.Confirm(fmt.Sprintf("Are you sure you want to delete the branch %s?", program.RenderElementSelected(branch.Name().Short())))

	if !confirmDelete {
		log.Info("Branch deletion cancelled.")
		return nil
	}

	refName := plumbing.ReferenceName(branch.Name())
	err := r.Storer.RemoveReference(refName)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Branch %s deleted successfully.", program.RenderElementSelected(branch.Name().Short())))

	return nil
}

func (r *repository) pull(branch *branch) error {
	if !r.isHead(branch) {
		log.Warn("You need to pull the branch from the HEAD, move on it first.")
		return nil
	}

	worktree, err := r.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err == git.NoErrAlreadyUpToDate {
		log.Info("Already up-to-date.")
		return nil
	}
	if err != nil {
		return err
	}

	log.Info("Branch pulled successfully.")

	return nil
}

func (r *repository) merge(branch *branch) error {
	if r.isHead(branch) {
		log.Warn("You cannot merge the branch you're currently on, please checkout on another branch first.")
		return nil
	}

	confirmMerge := program.Confirm(fmt.Sprintf("Are you sure you want to merge the branch %s into the current branch?", program.RenderElementSelected(branch.Name().Short())))

	if !confirmMerge {
		log.Info("Branch merge cancelled.")
		return nil
	}

	opts := &git.MergeOptions{
		Strategy: git.FastForwardMerge,
	}
	err := r.Merge(*branch, *opts)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Branch %s merged successfully into the current branch.", program.RenderElementSelected(branch.Name().Short())))

	return nil
}
