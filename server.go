package main

import (
	r "SDS/routes"
)

func main() {
	r.Routes.Run() // listen and serve on 0.0.0.0:8080
}
