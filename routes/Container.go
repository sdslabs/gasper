package routes

import (
	ContainerController "SDS/controllers/ContainerController"
)

func init() {
	Routes.POST("/container/create", ContainerController.Create)
}
