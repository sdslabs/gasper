package git

import (
	gogit "gopkg.in/src-d/go-git.v4"
)

// Clone clones the repo from the url into the destination
func Clone(url, destination string) error {
	_, err := gogit.PlainClone(destination, false, &gogit.CloneOptions{
		URL: url,
	})
	return err
}
