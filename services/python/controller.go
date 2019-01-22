package python

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/configs"
	g "github.com/sdslabs/SWS/lib/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/redis"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/lib/utils"
)

// createApp function handles requests for making making new static app
func createApp(c *gin.Context) {
	var (
		data map[string]interface{}
	)
	c.BindJSON(&data)

	ports, err := utils.GetFreePorts(2)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	if len(ports) < 2 {
		c.JSON(500, gin.H{
			"error": "Not Enough Ports",
		})
		return
	}

	sshPort := ports[0]
	httpPort := ports[1]

	context := data["context"].(map[string]interface{})

	appEnv, rer := api.CreateBasicApplication(
		data["name"].(string),
		data["url"].(string),
		strconv.Itoa(httpPort),
		strconv.Itoa(sshPort),
		context,
		&types.ApplicationConfig{
			DockerImage:  utils.ServiceConfig["python"].(map[string]interface{})["image"].(string),
			ConfFunction: configs.CreatePythonContainerConfig,
		})

	if rer != nil {
		g.SendResponse(c, rer, gin.H{})
		return
	}

	// Create virtual env for mentioned python
	_, rer = createVenv(appEnv, data["python_version"].(string))
	if rer != nil {
		g.SendResponse(c, rer, gin.H{})
		return
	}

	// Path of `requirements.txt` or any-other file containing requirements
	requirements := data["requirements"]
	if requirements != nil {
		_, rer = installRequirements(requirements.(string), appEnv)
		if rer != nil {
			g.SendResponse(c, rer, gin.H{})
			return
		}
	}

	// strore server process ID to kill it later on...
	execID, rer := startServer(data["run_command"].(string), appEnv)
	if rer != nil {
		g.SendResponse(c, rer, gin.H{})
		return
	}

	data["sshPort"] = sshPort
	data["httpPort"] = httpPort
	data["containerID"] = appEnv.ContainerID
	data["language"] = "python"
	data["hostIP"] = utils.HostIP
	data["execID"] = execID

	documentID, err := mongo.RegisterApp(data)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.RegisterApp(
		data["name"].(string),
		utils.HostIP+utils.ServiceConfig["python"].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	err = redis.IncrementServiceLoad(
		"python",
		utils.HostIP+utils.ServiceConfig["python"].(map[string]interface{})["port"].(string),
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"id":      documentID,
	})
}

func fetchDocs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "python"

	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

func deleteApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "python"

	c.JSON(200, gin.H{
		"message": mongo.DeleteApp(filter),
	})
}

func updateApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := utils.QueryToFilter(queries)

	filter["language"] = "python"

	var (
		data map[string]interface{}
	)
	c.BindJSON(&data)

	c.JSON(200, gin.H{
		"message": mongo.UpdateApp(filter, data),
	})
}
