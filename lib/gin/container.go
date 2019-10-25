package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// FetchLogs returns the container logs in a JSON format
func FetchLogs(c *gin.Context) {
	app := c.Param("app")
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	if filter["tail"] == nil {
		filter["tail"] = "-1"
	}

	data, err := docker.ReadLogs(app, filter["tail"].(string))

	if err != nil && err.Error() != "EOF" {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}

// FetchMysqlContainerLogs returns the mysql container logs in a JSON format
func FetchMysqlContainerLogs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	if filter["tail"] == nil {
		filter["tail"] = "-1"
	}

	data, err := docker.ReadLogs(types.MySQL, filter["tail"].(string))

	if err != nil && err.Error() != "EOF" {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}

// FetchMongoDBContainerLogs returns the mongodb container logs in a JSON format
func FetchMongoDBContainerLogs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	if filter["tail"] == nil {
		filter["tail"] = "-1"
	}

	data, err := docker.ReadLogs(types.MongoDB, filter["tail"].(string))

	if err != nil && err.Error() != "EOF" {
		utils.SendServerErrorResponse(c, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}

// ReloadServer reloads the nginx server
func ReloadServer(c *gin.Context) {
	app := c.Param("app")

	cmd := []string{"nginx", "-s", "reload"}
	_, err := docker.ExecDetachedProcess(app, cmd)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

// ReloadMysqlService reloads the Mysql service in the container
func ReloadMysqlService(c *gin.Context) {
	cmd := []string{"service", "mysql", "start"}
	_, err := docker.ExecDetachedProcess(types.MySQL, cmd)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

// ReloadMongoDBService reloads the Mysql service in the container
func ReloadMongoDBService(c *gin.Context) {
	cmd := []string{"service", "monogdb", "restart"}
	_, err := docker.ExecDetachedProcess(types.MongoDB, cmd)
	if err != nil {
		utils.SendServerErrorResponse(c, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}
