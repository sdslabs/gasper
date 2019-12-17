package middlewares

import (
	"encoding/json"
	"fmt"

	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

var disallowedApplicationNames = []string{
	types.Kaze,
	types.Mizu,
	types.Hikari,
	types.Enrai,
	types.MySQL,
	types.MongoDB,
	types.Kaen,
	types.Iwa,
}

var disallowedDatabaseNames = []string{
	"admin",
	"config",
	"local",
	"root",
	"mysql",
	"information_schema",
	"performance_schema",
	"sys",
}

func isUniqueInstance(instanceName, instanceType string) (bool, error) {
	count, err := mongo.CountInstances(types.M{
		mongo.NameKey:         instanceName,
		mongo.InstanceTypeKey: instanceType,
	})
	if err != nil || count != 0 {
		return false, err
	}
	return true, nil
}

// ValidateApplicationRequest validates the request for creating applications
func ValidateApplicationRequest(c *gin.Context) {
	requestBody := getBodyFromContext(c)
	app := &types.ApplicationConfig{}
	err := json.Unmarshal(requestBody, app)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if result, err := validator.ValidateStruct(app); !result {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if utils.Contains(disallowedApplicationNames, app.GetName()) {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Name of application cannot be `%s`", app.GetName()),
		})
		return
	}

	unique, err := isUniqueInstance(app.GetName(), mongo.AppInstance)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	if !unique {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Application with that name already exists",
		})
		return
	}
	c.Next()
}

// ValidateDatabaseRequest validates the request for creating databases
func ValidateDatabaseRequest(c *gin.Context) {
	requestBody := getBodyFromContext(c)
	db := &types.DatabaseConfig{}
	err := json.Unmarshal(requestBody, db)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if result, err := validator.ValidateStruct(db); !result {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if utils.Contains(disallowedDatabaseNames, db.GetName()) {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Name of database cannot be `%s`", db.GetName()),
		})
		return
	}

	unique, err := isUniqueInstance(db.GetName(), mongo.DBInstance)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}

	if !unique {
		c.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Database with that name already exists",
		})
		return
	}
	c.Next()
}
