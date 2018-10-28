package routes

import (
	phpAppController "github.com/sdslabs/SDS/controllers/phpapp"
)

func init() {
	appGroup := Router.Group("/app/php")
	{
		appGroup.POST("/create", phpAppController.CreateApp)
	}
}
