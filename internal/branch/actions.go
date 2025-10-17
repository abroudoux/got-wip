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
		return delete(branchSelected, r)
	case actionNewBranch:
		return createNewBranch(branchSelected, r)
	case actionCopyName:
		return copyBranchName(branchSelected)
	case actionCheckout:
		return checkout(branchSelected, r)
	case actionPull:
		if !isHead(branchSelected, r) {
			log.Warn("You need to pull the branch from the HEAD, move on it first.")
			return nil
		}

		return pull(r)
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

func checkout(branch *branch, r *repository) error {
	if isHead(branch, r) {
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

func createNewBranch(branch *branch, r *repository) error {
	if !isHead(branch, r) {
		log.Warn("You need to create a branch from the HEAD, move on it first.")
		return nil
	}

	var newBranchName string
	var checkoutOnNewBranch bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter the name of the new branch:").Value(&newBranchName).Validate(
				func(input string) error {
					if isBranchNameAlreadyExists(input, r) {
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

	head, err := getHead(r)
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

	if checkoutOnNewBranch {
		branchCreated, err := findBranchByName(newBranch.Name().Short(), r)
		if err != nil {
			return err
		}
		if branchCreated == nil {
			return fmt.Errorf("branch %s not found after creation", newBranch.Name().Short())
		}

		err = checkout(branchCreated, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func delete(branch *branch, r *repository) error {
	if isHead(branch, r) {
		log.Warn("You cannot delete the branch you're currently on, please checkout on another branch first.")
		return nil
	}

	var confirmDelete bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Are you sure you want to delete " + branch.Name().Short() + "?").Value(&confirmDelete),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	if !confirmDelete {
		log.Info("Branch deletion cancelled.")
		return nil
	}

	refName := plumbing.ReferenceName(branch.Name())
	err = r.Storer.RemoveReference(refName)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Branch %s deleted successfully.", program.RenderElementSelected(branch.Name().Short())))

	return nil
}

func pull(r *repository) error {
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
