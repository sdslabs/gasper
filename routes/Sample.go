package routes

import (
	SampleController "SDS/controllers/SampleController"
)

func init() {
	Routes.GET("/ping", SampleController.Pong)
}
