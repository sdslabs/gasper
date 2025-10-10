package types

type UpNodeRequest struct {
	//ip:port of the node
	NodeIP string `json:"node_ip" binding:"required"`
}

type Response struct {
	Message string `json:"message"`
}
