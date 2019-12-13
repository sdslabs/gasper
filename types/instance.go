package types

// Instance is the interface for dealing with both applications and databases
type Instance interface {
	GetName() string
}
