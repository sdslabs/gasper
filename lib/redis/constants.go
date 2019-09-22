package redis

// AppKey is the key name for the HashMap containing application instances
const AppKey string = "apps"

// DatabaseKey is the key name for the HashMap containing database instances
const DatabaseKey string = "databases"

// SSHKey is the key name for the Sorted Set containing ssh microservice instances
const SSHKey string = "ssh"

// ErrEmptySet is the error message when the redis set being queried is empty
const ErrEmptySet string = "Empty Set"
