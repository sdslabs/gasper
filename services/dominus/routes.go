package dominus

import (
	"net/http"
	"time"

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

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cookie"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig), middlewares.FalconGuard())

	auth := router.Group("/auth")
	{
		auth.POST("/login", middlewares.JWT.LoginHandler)
		auth.POST("/register", middlewares.RegisterValidator, middlewares.Register)
		auth.GET("/refresh", middlewares.JWT.MiddlewareFunc(), middlewares.JWT.RefreshHandler)
	}

	app := router.Group("/apps")
	app.Use(middlewares.JWT.MiddlewareFunc())
	{
		app.POST("/:language", middlewares.ValidateApplicationRequest, createApp)
		app.GET("", fetchAppsByUser())
		app.GET("/:app", middlewares.IsAppOwner(), gin.FetchAppInfo)
		app.PUT("/:app", middlewares.IsAppOwner(), gin.UpdateAppByName)
		app.DELETE("/:app", middlewares.IsAppOwner(), deleteApp)
		app.GET("/:app/logs", middlewares.IsAppOwner(), fetchAppLogs)
		app.PATCH("/:app/rebuild", middlewares.IsAppOwner(), rebuildApp)
	}

	db := router.Group("/dbs")
	db.Use(middlewares.JWT.MiddlewareFunc())
	{
		db.POST("/:database", middlewares.ValidateDatabaseRequest, trimURLPath(2), createDatabase)
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
			apps.DELETE("/:app", deleteApp)
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
			users.PATCH("/:user/grant", adminHandlers.GrantSuperuserPrivilege())
			users.PATCH("/:user/revoke", adminHandlers.RevokeSuperuserPrivilege())
		}
		nodes := admin.Group("/nodes")
		{
			nodes.GET("", adminHandlers.GetAllNodes)
			nodes.GET("/:type", adminHandlers.GetNodesByName)
		}
	}

	return router
}
