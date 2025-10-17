package branch

func isHead(branch *branch, r *repository) bool {
	head, err := getHead(r)
	if err != nil {
		return false
	}

	return head.Name().Short() == branch.Name().Short()
}

func isBranchNameAlreadyExists(branchName string, r *repository) bool {
	branches, err := getBranches(r)
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
