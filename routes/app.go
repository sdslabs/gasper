package routes

import (
	appController "github.com/sdslabs/SDS/controllers/app"
)

func init() {
	appGroup := Router.Group("/app")
	{
		createGroup := appGroup.Group("/create")
		{
			createGroup.POST("/static", appController.CreateStaticApp)
		}
	}
}
