package types

type UpNodeRequest struct {
	//ip:port of the node
	NodeAddress string `json:"node_address" binding:"required"`
}

type Response struct {
	Message string `json:"message"`
}
