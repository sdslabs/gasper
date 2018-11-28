package main

import (
	r "github.com/sdslabs/SWS/routes"
)

func main() {
	// listen and serve on 0.0.0.0:8080
	r.Router.Run()
}
