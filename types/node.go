package types

type NodeRequest struct {
	Node         string `json:"node" validate:"required,hostname_port"`
	DeleteVolume bool   `json:"delete_volumes"`
}
