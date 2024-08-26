package master

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	c "github.com/sdslabs/gasper/services/master/controllers"
	m "github.com/sdslabs/gasper/services/master/middlewares"
	"github.com/sdslabs/gasper/types"
)

// ServiceName is the name of the current microservice
const ServiceName = types.Master

// NewService returns a new instance of the current microservice
func NewService() http.Handler {
	// router is the main routes handler for the current microservice package
	router := gin.Default()

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cookie", "Authorization-Type"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))
	router.NoRoute(c.Handle404)

	// Bind frontend generated from https://github.com/sdslabs/SWS
	router.GET("", func(ctx *gin.Context) {
		ctx.Data(200, frontendBinder["index.html"].responseHeader, frontendBinder["index.html"].content)
	})
	for file, box := range frontendBinder {
		// A deep copy of the filename and the box is made because as the loop iterator traverses
		// all of the handler functions point to the last element of the map which is not correct
		// This is due to the iterator occupying a single instance of heap memory and all handler
		// functions pointing to that single heap memory instance
		// Making a deep copy makes separate clones of heap memory instances thereby preventing
		// the undesired override
		fileDeepCopy := file
		boxDeepCopy := box
		router.GET("/"+fileDeepCopy, func(ctx *gin.Context) {
			ctx.Data(200, boxDeepCopy.responseHeader, boxDeepCopy.content)
		})
	}

	auth := router.Group("/auth")
	{
		auth.POST("/login", m.LoginHandler)
		auth.POST("/register", m.ValidateRegistration, c.Register)
		auth.GET("/refresh", m.RefreshHandler)
		auth.PUT("/revoke", c.RevokeToken)
	}

	router.GET("/instances", m.AuthRequired(), c.FetchAllInstancesByUser)
	router.POST("/gctllogin", m.JWTGctl.MiddlewareFunc(), c.GctlLogin)
	router.POST("/github", m.AuthRequired(), c.CreateRepository)

	app := router.Group("/apps")
	app.Use(m.AuthRequired())
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
		app.GET("/:app/metrics", m.IsAppOwner, c.FetchMetrics)
	}

	db := router.Group("/dbs")
	db.Use(m.AuthRequired())
	{
		db.POST("/:database", m.ValidateDatabaseRequest, c.CreateDatabase)
		db.GET("", c.FetchDatabasesByUser)
		db.GET("/:db", m.IsDatabaseOwner, c.GetDatabaseInfo)
		db.DELETE("/:db", m.IsDatabaseOwner, c.DeleteDatabase)
		db.PATCH("/:db/transfer/:user", m.IsDatabaseOwner, c.TransferDatabaseOwnership)
		db.GET("/:db/redislogs",m.IsDatabaseOwner,c.GetRedisLogs)
	}

	user := router.Group("/user")
	user.Use(m.AuthRequired())
	{
		user.GET("", c.GetLoggedInUserInfo)
		user.PUT("/password", c.UpdatePassword)
		user.DELETE("", c.DeleteUser)
	}

	admin := router.Group("/admin")
	admin.Use(m.AuthRequired(), m.VerifyAdmin)
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
