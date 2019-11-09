package git

import (
	"fmt"

	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Clone clones the repo from the url into the destination
func Clone(url, branch, destination string) error {
	_, err := gogit.PlainClone(destination, false, &gogit.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Depth:         1,
	})
	return err
}

// CloneWithToken clones the repo from the url using access token
// which is an alternative for auth to clone private repositories
func CloneWithToken(url, branch, destination, token string) error {
	_, err := gogit.PlainClone(destination, false, &gogit.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Depth:         1,
		Auth: &http.BasicAuth{
			// Since cloning through token username can be anything but not an empty string
			Username: "random_string",
			Password: token,
		},
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
