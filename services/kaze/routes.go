package kaze

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	c "github.com/sdslabs/gasper/services/kaze/controllers"
	m "github.com/sdslabs/gasper/services/kaze/middlewares"
	"github.com/sdslabs/gasper/types"
)

// ServiceName is the name of the current microservice
const ServiceName = types.Kaze

// NewService returns a new instance of the current microservice
func NewService() http.Handler {
	// router is the main routes handler for the current microservice package
	router := gin.Default()

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cookie"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig), m.FalconGuard())
	router.NoRoute(c.Handle404)

	auth := router.Group("/auth")
	{
		auth.POST("/login", m.Register, m.JWT.LoginHandler)
		auth.GET("/refresh", m.JWT.RefreshHandler)
	}

	app := router.Group("/apps")
	app.Use(m.JWT.MiddlewareFunc())
	{
		app.POST("/:language", m.ValidateApplicationRequest, c.CreateApp)
		app.GET("", c.FetchAppsByUser)
		app.GET("/:app", m.IsAppOwner, c.GetApplicationInfo)
		app.PUT("/:app", m.IsAppOwner, c.UpdateAppByName)
		app.DELETE("/:app", m.IsAppOwner, c.DeleteApp)
		app.GET("/:app/logs", m.IsAppOwner, c.FetchAppLogs)
		app.PATCH("/:app/rebuild", m.IsAppOwner, c.RebuildApp)
		app.PATCH("/:app/transfer/:user", m.IsAppOwner, c.TransferApplicationOwnership)
		app.GET("/:app/term", m.IsAppOwner, c.DeployWebTerminal)
		app.GET("/:app/metrics", c.FetchMetrics)
	}

	db := router.Group("/dbs")
	db.Use(m.JWT.MiddlewareFunc())
	{
		db.POST("/:database", m.ValidateDatabaseRequest, c.CreateDatabase)
		db.GET("", c.FetchDatabasesByUser)
		db.GET("/:db", m.IsDatabaseOwner, c.GetDatabaseInfo)
		db.DELETE("/:db", m.IsDatabaseOwner, c.DeleteDatabase)
		db.PATCH("/:db/transfer/:user", m.IsDatabaseOwner, c.TransferDatabaseOwnership)
	}

	user := router.Group("/user")
	user.Use(m.JWT.MiddlewareFunc())
	{
		user.GET("", c.GetLoggedInUserInfo)
		user.PUT("/password", c.UpdatePassword)
		user.DELETE("", c.DeleteUser)
	}

	admin := router.Group("/admin")
	admin.Use(m.JWT.MiddlewareFunc(), m.VerifyAdmin)
	{
		apps := admin.Group("/apps")
		{
			apps.GET("", c.GetAllApplications)
			apps.GET("/:app", c.GetApplicationInfo)
			apps.DELETE("/:app", c.DeleteApp)
		}
		dbs := admin.Group("/dbs")
		{
			dbs.GET("", c.GetAllDatabases)
			dbs.GET("/:db", c.GetDatabaseInfo)
			dbs.DELETE("/:db", c.DeleteDatabase)
		}
		users := admin.Group("/users")
		{
			users.GET("", c.GetAllUsers)
			users.GET("/:user", c.GetUserInfo)
			users.DELETE("/:user", c.DeleteUserByAdmin)
			users.PATCH("/:user/grant", c.GrantSuperuserPrivilege)
			users.PATCH("/:user/revoke", c.RevokeSuperuserPrivilege)
		}
		nodes := admin.Group("/nodes")
		{
			nodes.GET("", c.GetAllNodes)
			nodes.GET("/:type", c.GetNodesByName)
		}
	}

	return router
}
