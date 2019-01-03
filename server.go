package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
	"github.com/sdslabs/SWS/services/dominus"
	"github.com/sdslabs/SWS/services/php"
	"github.com/sdslabs/SWS/services/static"
	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group

	// Bind services to routers here
	serviceBindings := map[string]*gin.Engine{
		"dominus": dominus.Router,
		"static":  static.Router,
		"php":     php.Router,
	}

	for service, config := range utils.ServiceConfig {
		config := config.(map[string]interface{})
		if config["deploy"].(bool) {
			fmt.Printf("%s Service Active\n", strings.Title(service))
			server := &http.Server{
				Addr:         config["port"].(string),
				Handler:      serviceBindings[service],
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 30 * time.Second,
			}
			g.Go(func() error {
				return server.ListenAndServe()
			})
		}
	}

	dominus.ExposeServices()

	if utils.ServiceConfig["dominus"].(map[string]interface{})["deploy"].(bool) {
		cleanupInterval := time.Duration(utils.SWSConfig["cleanupInterval"].(float64))
		dominus.ScheduleCleanup(cleanupInterval * time.Second)
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
