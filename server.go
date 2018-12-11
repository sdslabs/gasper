package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
	"github.com/sdslabs/SWS/services/static"
	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group

	// Bind services to routers here
	serviceBindings := map[string]*gin.Engine{
		"static": static.Router,
	}

	for _, service := range utils.SWSConfig.Services {
		if service.Deploy {
			for _, port := range service.Ports {
				server := &http.Server{
					Addr:         port,
					Handler:      serviceBindings[service.Name],
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}
				g.Go(func() error {
					return server.ListenAndServe()
				})
			}
		}
	}

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
