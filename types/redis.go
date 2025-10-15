package types

// InstanceBindings defines the struct for storing both the instance's server and node urls
type InstanceBindings struct {
	Node   string `json:"node" valid:"required~Field 'node' is required,matches(^([0-9]|[0-9][0-9]|[0-9][0-9][0-9]\\.){3}([0-9]|[0-9][0-9]|[0-9][0-9][0-9]):([0-9]|[0-9][0-9]|[0-9][0-9][0-9][0-9]|[0-9][0-9][0-9][0-9][0-9])$)~Field 'node' must match format 'x.x.x.x:PORT' (appropriate PORT number)"`
	Server string `json:"server"`
}
