package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/utils"
	"github.com/sdslabs/SWS/services/php"
	"github.com/sdslabs/SWS/services/static"
	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group

	// Bind services to routers here
	serviceBindings := map[string]*gin.Engine{
		"static": static.Router,
		"php":    php.Router,
	}

	for _, service := range utils.SWSConfig.Services {
		if service.Deploy {
			server := &http.Server{
				Addr:         service.Port,
				Handler:      serviceBindings[service.Name],
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}
			g.Go(func() error {
				return server.ListenAndServe()
			})
		}
	}

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
