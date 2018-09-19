package routes

import (
	containerController "github.com/sdslabs/SDS/controllers/container"
)

func init() {
	containerGroup := Router.Group("/container")
	{
		containerGroup.GET("/create", containerController.Create)
	}
}
