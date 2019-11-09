package redis

import "github.com/sdslabs/gasper/types"

const (
	// ApplicationKey is the key name for the HashMap containing application instances
	ApplicationKey string = "applications"

	// DatabaseKey is the key name for the HashMap containing database instances
	DatabaseKey string = "databases"

	// SSHKey is the key name for the Sorted Set containing ssh microservice instances
	SSHKey string = types.Iwa

	// WorkerInstanceKey is the key name for Worker nodes
	WorkerInstanceKey string = types.Mizu

	// ErrEmptySet is the error message when the redis set being queried is empty
	ErrEmptySet string = "Empty Set"
)
