package routes

import (
	staticAppController "github.com/sdslabs/SWS/controllers/staticapp"
)

func init() {
	appGroup := Router.Group("/app/static")
	{
		appGroup.POST("/create", staticAppController.CreateApp)
	}
}
