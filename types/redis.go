package types

// AppBindings defines the struct for storing both the server and node urls
type AppBindings struct {
	Node   string `json:"node"`
	Server string `json:"server"`
}
