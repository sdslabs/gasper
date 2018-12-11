package static

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/mongo"
	"github.com/sdslabs/SWS/lib/utils"
)

// staticAppConfig is json binding config for creating new static page
type staticAppConfig struct {
	Name      string `json:"name" form:"name" binding:"required"`
	UserID    int    `json:"user_id" form:"user_id" binding:"required"`
	GithubURL string `json:"github_url" form:"github_url" binding:"required"`
}

// readAndWriteConfig creates new config file for the given app
func ReadAndWriteConfig(name string) utils.Error {
	// containerID, ok := os.LookupEnv("STATIC_CONTAINER_ID")
	// if !ok {
	// 	return utils.Error{
	// 		Code: 500,
	// 		Err:  errors.New("STATIC_CONTAINER_ID not found in the environment"),
	// 	}
	// }

	err := utils.ReadAndWriteConfig(name, "static", "3b99fa7534c3")
	if err != nil {
		return utils.Error{
			Code: 500,
			Err:  err,
		}
	}

	return utils.Error{
		Code: 200,
		Err:  nil,
	}
}

// createApp function handles requests for making making new static app
func createApp(c *gin.Context) {
	var (
		data map[string]interface{}
		// err  utils.Error
	)
	c.BindJSON(&data)
	data["language"] = "static"

	// err = ReadAndWriteConfig("random")
	// if err.Code != 200 {
	// 	c.JSON(err.Code, gin.H{
	// 		"message": err.Reason(),
	// 	})
	// 	return
	// }

	c.JSON(200, gin.H{
		"success": true,
		"id":      mongo.RegisterApp(data),
	})
}

func fetchDocs(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := queryToFilter(queries)

	filter["language"] = "static"

	c.JSON(200, gin.H{
		"data": mongo.FetchAppInfo(filter),
	})
}

func deleteApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := queryToFilter(queries)

	filter["language"] = "static"

	c.JSON(200, gin.H{
		"message": mongo.DeleteApp(filter),
	})
}

func updateApp(c *gin.Context) {
	queries := c.Request.URL.Query()
	filter := queryToFilter(queries)

	filter["language"] = "static"

	var (
		data map[string]interface{}
		// err  utils.Error
	)
	c.BindJSON(&data)

	c.JSON(200, gin.H{
		"message": mongo.UpdateApp(filter, data),
	})
}

func queryToFilter(queries map[string][]string) map[string]interface{} {
	filter := make(map[string]interface{})

	for key, value := range queries {
		filter[key] = value[0]
	}

	return filter
}
