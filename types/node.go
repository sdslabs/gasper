package types

type NodeRequest struct {
	Node         string `json:"node" valid:"required,dialstring"`
	DeleteVolume bool   `json:"delete_volumes"`
}
