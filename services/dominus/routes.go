package dominus

import (
	"github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
	adminHandlers "github.com/sdslabs/SWS/services/dominus/admin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewEngine()

// ServiceName is the name of the current microservice
var ServiceName = "dominus"

func init() {
	Router.Use(middlewares.FalconGuard())
	Router.POST("/:service", createApp)
	// Router.GET("/", gin.FetchDocs(ServiceName))
	// Router.PUT("/", gin.UpdateAppInfo(ServiceName))
	// Router.DELETE("/", gin.DeleteApp(ServiceName))
	app := Router.Group("/apps")
	{
		app.GET("/:app", gin.FetchAppInfo)
		app.PUT("/:app", gin.UpdateAppByName)
		app.DELETE("/:app", gin.DeleteAppByName)
		app.GET("/:app/:action", trimURLPath, execute)
	}
	db := Router.Group("/dbs")
	{
		db.GET("/:db", gin.FetchDBInfo)
		db.DELETE("/:user/:db", trimURLPath, deleteDB)
	}
	admin := Router.Group("/admin")
	{
		apps := admin.Group("/apps")
		{
			apps.GET("", adminHandlers.GetAllApplications)
			apps.GET("/:app", adminHandlers.GetApplicationInfo)
			apps.DELETE("/:app", adminHandlers.DeleteApplication)
		}
		dbs := admin.Group("/dbs")
		{
			dbs.GET("", adminHandlers.GetAllDatabases)
			dbs.GET("/:db", adminHandlers.GetDatabaseInfo)
			dbs.DELETE("/:user/:db", trimURLPath, execute)
		}
		users := admin.Group("/users")
		{
			users.GET("", adminHandlers.GetAllUsers)
			users.GET("/:user", adminHandlers.GetUserInfo)
			users.DELETE("/:user", adminHandlers.DeleteUser)
		}
		nodes := admin.Group("/nodes")
		{
			nodes.GET("", adminHandlers.GetAllNodes)
			nodes.GET("/:node", adminHandlers.GetNodesByName)
		}
	}
}
