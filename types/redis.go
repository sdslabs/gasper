package types

// InstanceBindings defines the struct for storing both the instance's server and node urls
type InstanceBindings struct {
	Node   string `json:"node"`
	Server string `json:"server"`
}
