package dominus

import (
	"github.com/gin-contrib/cors"
	"github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/middlewares"
	adminHandlers "github.com/sdslabs/SWS/services/dominus/admin"
)

// Router is the main routes handler for the current microservice package
var Router = gin.NewEngine()

// ServiceName is the name of the current microservice
var ServiceName = "dominus"

func init() {
	Router.Use(cors.Default(), middlewares.FalconGuard())
	auth := Router.Group("/auth")
	{
		auth.POST("/login", middlewares.JWT.LoginHandler)
		auth.POST("/register", middlewares.RegisterValidator, middlewares.Register)
		auth.GET("/refresh", middlewares.JWT.MiddlewareFunc(), middlewares.JWT.RefreshHandler)
	}
	app := Router.Group("/apps")
	app.Use(middlewares.JWT.MiddlewareFunc())
	{
		app.POST("/:language", trimURLPath(2), createApp)
		app.GET("/:app", gin.FetchAppInfo)
		app.PUT("/:app", gin.UpdateAppByName)
		app.DELETE("/:app", trimURLPath(2), execute)
		app.GET("/:app/:action", trimURLPath(2), execute)
	}
	db := Router.Group("/dbs")
	db.Use(middlewares.JWT.MiddlewareFunc())
	{
		db.POST("/:database", trimURLPath(2), createDatabase)
		db.GET("/:db", gin.FetchDBInfo)
		db.DELETE("/:db", trimURLPath(2), deleteDB)
	}
	admin := Router.Group("/admin")
	admin.Use(middlewares.JWT.MiddlewareFunc(), middlewares.VerifyAdmin)
	{
		apps := admin.Group("/apps")
		{
			apps.GET("/", adminHandlers.GetAllApplications)
			apps.GET("/:app", adminHandlers.GetApplicationInfo)
			apps.DELETE("/:app", trimURLPath(3), execute)
		}
		dbs := admin.Group("/dbs")
		{
			dbs.GET("/", adminHandlers.GetAllDatabases)
			dbs.GET("/:db", adminHandlers.GetDatabaseInfo)
			dbs.DELETE("/:db", trimURLPath(3), deleteDB)
		}
		users := admin.Group("/users")
		{
			users.GET("/", adminHandlers.GetAllUsers)
			users.GET("/:user", adminHandlers.GetUserInfo)
			users.DELETE("/:user", adminHandlers.DeleteUser)
		}
		nodes := admin.Group("/nodes")
		{
			nodes.GET("/", adminHandlers.GetAllNodes)
			nodes.GET("/:node", adminHandlers.GetNodesByName)
		}
	}
}
