package routes

import (
	"github.com/sdslabs/SDS/controllers/app"
)

func init() {
	appGroup := Router.Group("/app")
	{
		appGroup.POST("/create", app.Create)
	}
}
