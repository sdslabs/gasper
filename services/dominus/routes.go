package dominus

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/sdslabs/gasper/lib/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
	adminHandlers "github.com/sdslabs/gasper/services/dominus/admin"
	"github.com/sdslabs/gasper/types"
)

// ServiceName is the name of the current microservice
const ServiceName = types.Dominus

// NewService returns a new instance of the current microservice
func NewService() http.Handler {
	// router is the main routes handler for the current microservice package
	router := gin.NewEngine()
	router.Use(cors.Default(), middlewares.FalconGuard())

	auth := router.Group("/auth")
	{
		auth.POST("/login", middlewares.JWT.LoginHandler)
		auth.POST("/register", middlewares.RegisterValidator, middlewares.Register)
		auth.GET("/refresh", middlewares.JWT.MiddlewareFunc(), middlewares.JWT.RefreshHandler)
	}

	app := router.Group("/apps")
	app.Use(middlewares.JWT.MiddlewareFunc())
	{
		app.POST("/:language", middlewares.InsertOwner, trimURLPath(2), createApp)
		app.GET("", fetchAppsByUser())
		app.GET("/:app", middlewares.IsAppOwner(), gin.FetchAppInfo)
		app.PUT("/:app", middlewares.IsAppOwner(), gin.UpdateAppByName)
		app.DELETE("/:app", middlewares.IsAppOwner(), trimURLPath(2), execute)
		app.GET("/:app/:action", middlewares.IsAppOwner(), trimURLPath(2), execute)
	}

	db := router.Group("/dbs")
	db.Use(middlewares.JWT.MiddlewareFunc())
	{
		db.POST("/:database", middlewares.InsertOwner, trimURLPath(2), createDatabase)
		db.GET("", fetchDBsByUser())
		db.GET("/:db", middlewares.IsDbOwner(), gin.FetchDBInfo)
		db.DELETE("/:db", middlewares.IsDbOwner(), trimURLPath(2), deleteDB)
	}

	admin := router.Group("/admin")
	admin.Use(middlewares.JWT.MiddlewareFunc(), middlewares.VerifyAdmin)
	{
		apps := admin.Group("/apps")
		{
			apps.GET("", adminHandlers.GetAllApplications)
			apps.GET("/:app", adminHandlers.GetApplicationInfo)
			apps.DELETE("/:app", trimURLPath(3), execute)
		}
		dbs := admin.Group("/dbs")
		{
			dbs.GET("", adminHandlers.GetAllDatabases)
			dbs.GET("/:db", adminHandlers.GetDatabaseInfo)
			dbs.DELETE("/:db", trimURLPath(3), deleteDB)
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

	return router
}
