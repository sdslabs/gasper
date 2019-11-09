package dominus

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/middlewares"
	. "github.com/sdslabs/gasper/services/dominus/controllers"
	"github.com/sdslabs/gasper/types"
)

// ServiceName is the name of the current microservice
const ServiceName = types.Dominus

// NewService returns a new instance of the current microservice
func NewService() http.Handler {
	// router is the main routes handler for the current microservice package
	router := gin.Default()

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
		app.POST("/:language", middlewares.ValidateApplicationRequest, CreateApp)
		app.GET("", FetchAppsByUser)
		app.GET("/:app", middlewares.IsAppOwner(), GetApplicationInfo)
		app.PUT("/:app", middlewares.IsAppOwner(), UpdateAppByName)
		app.DELETE("/:app", middlewares.IsAppOwner(), DeleteApp)
		app.GET("/:app/logs", middlewares.IsAppOwner(), FetchAppLogs)
		app.PATCH("/:app/rebuild", middlewares.IsAppOwner(), RebuildApp)
	}

	db := router.Group("/dbs")
	db.Use(middlewares.JWT.MiddlewareFunc())
	{
		db.POST("/:database", middlewares.ValidateDatabaseRequest, CreateDatabase)
		db.GET("", FetchDatabasesByUser)
		db.GET("/:db", middlewares.IsDbOwner(), GetDatabaseInfo)
		db.DELETE("/:db", middlewares.IsDbOwner(), DeleteDatabase)
	}

	admin := router.Group("/admin")
	admin.Use(middlewares.JWT.MiddlewareFunc(), middlewares.VerifyAdmin)
	{
		apps := admin.Group("/apps")
		{
			apps.GET("", GetAllApplications)
			apps.GET("/:app", GetApplicationInfo)
			apps.DELETE("/:app", DeleteApp)
		}
		dbs := admin.Group("/dbs")
		{
			dbs.GET("", GetAllDatabases)
			dbs.GET("/:db", GetDatabaseInfo)
			dbs.DELETE("/:db", DeleteDatabase)
		}
		users := admin.Group("/users")
		{
			users.GET("", GetAllUsers)
			users.GET("/:user", GetUserInfo)
			users.DELETE("/:user", DeleteUser)
			users.PATCH("/:user/grant", GrantSuperuserPrivilege)
			users.PATCH("/:user/revoke", RevokeSuperuserPrivilege)
		}
		nodes := admin.Group("/nodes")
		{
			nodes.GET("", GetAllNodes)
			nodes.GET("/:type", GetNodesByName)
		}
	}

	return router
}
