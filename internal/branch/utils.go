package branch

func (r *repository) isHead(branch *branch) bool {
	return r.head.Name().Short() == branch.Name().Short()
}

func (r *repository) isBranchNameAlreadyExists(branchName string) bool {
	for _, b := range r.branches {
		if b.Name().Short() == branchName {
			return true
		}
	}

	return false
}
