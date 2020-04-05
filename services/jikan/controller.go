package jikan

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

// ServiceName is the name of the current microservice
const ServiceName = types.Jikan

func streamHandler(c *gin.Context) {
	appName := c.Param("app")

	metricsInterval, err := strconv.ParseInt(c.Query("interval"), 10, 64)
	if err != nil {
		metricsInterval = 2 * int64(configs.ServiceConfig.Mizu.MetricsInterval)
	}

	metricsCount, err := strconv.ParseInt(c.Query("count"), 10, 64)
	if err != nil {
		metricsCount = 10
	}

	chanStream := make(chan []types.M, 10)
	go func() {
		for {
			defer close(chanStream)
			metrics := mongo.FetchContainerMetrics(types.M{
				mongo.NameKey: appName,
				mongo.TimestampKey: types.M{
					"$gte": time.Now().Unix() - int64(configs.ServiceConfig.Mizu.MetricsInterval*time.Second),
				},
			}, metricsCount)
			chanStream <- metrics
			if metricsInterval < int64(configs.ServiceConfig.Mizu.MetricsInterval) {
				metricsInterval = 2 * int64(configs.ServiceConfig.Mizu.MetricsInterval)
			}

			time.Sleep(time.Second * time.Duration(metricsInterval))
		}
	}()
	c.Stream(func(w io.Writer) bool {
		if metrics, ok := <-chanStream; ok {
			c.SSEvent("metrics", metrics)
			return true
		}

		return false
	})
}

// NewService returns a new instance of the current microservice
func NewService() *http.Server {
	if !utils.IsValidPort(configs.ServiceConfig.Jikan.Port) {
		msg := fmt.Sprintf("Port %d is invalid or already in use.\n", configs.ServiceConfig.Jikan.Port)
		utils.Log(msg, utils.ErrorTAG)
		os.Exit(1)
	}

	router := gin.Default()

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET"},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           72 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	router.GET("/stream/:app/metrics", streamHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.ServiceConfig.Jikan.Port),
		Handler: router,
	}
	return server
}
