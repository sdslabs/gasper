package types

// InstanceBindings defines the struct for storing both the instance's server and node urls
type InstanceBindings struct {
	Node   string `json:"node" validate:"required,hostname_port"`
	Server string `json:"server"`
}
