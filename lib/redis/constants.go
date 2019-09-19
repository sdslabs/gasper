package redis

const (
	// AppKey is the key name for the HashMap containing application instances
	AppKey string = "apps"

	// DatabaseKey is the key name for the HashMap containing database instances
	DatabaseKey string = "databases"

	// SSHKey is the key name for the Sorted Set containing ssh microservice instances
	SSHKey string = "ssh"

	// ErrEmptySet is the error message when the redis set being queried is empty
	ErrEmptySet string = "Empty Set"
)
