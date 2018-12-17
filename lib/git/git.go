package git

import (
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Clone clones the repo from the url into the destination
func Clone(url, destination string) error {
	_, err := gogit.PlainClone(destination, false, &gogit.CloneOptions{
		URL: url,
	})
	return err
}

// Pull pulls the latest branch from "origin"
// 'dotgitPath' is the absolue path to .git directory
// 'branch' is the branch name which is to be pulled
func Pull(dotgitPath, branch string) error {
	repo, err := gogit.PlainOpen(dotgitPath)
	if err != nil {
		return err
	}
	wtree, err := repo.Worktree()
	if err != nil {
		return err
	}
	refName := plumbing.NewBranchReferenceName(branch)
	pullOpts := &gogit.PullOptions{
		RemoteName:    "origin",
		ReferenceName: refName,
		SingleBranch:  true,
	}
	err = wtree.Pull(pullOpts)
	if err == gogit.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}
