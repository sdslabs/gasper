package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	static "github.com/sdslabs/SWS/services/static"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	Name   string `json:"name"`
	Deploy bool   `json:"deploy"`
	Port   string `json:"port"`
}

type Services struct {
	Services []Service `json:"services"`
}

func main() {

	var g errgroup.Group

	file, _ := ioutil.ReadFile("./config.json")
	var services Services
	err := json.Unmarshal(file, &services)

	if err != nil {
		panic(err)
	}

	// Bind services to routers here
	serviceBindings := map[string]*gin.Engine{
		"static": static.Router,
	}

	for _, service := range services.Services {

		if service.Deploy == true {

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
