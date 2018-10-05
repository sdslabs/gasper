package routes

import (
	staticAppController "github.com/sdslabs/SDS/controllers/staticapp"
)

func init() {
	appGroup := Router.Group("/app/static")
	{
		appGroup.POST("/create", staticAppController.CreateApp)
	}
}
