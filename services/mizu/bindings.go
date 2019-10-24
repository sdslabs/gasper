package mizu

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/services/mizu/components/nodejs"
	"github.com/sdslabs/gasper/services/mizu/components/php"
	"github.com/sdslabs/gasper/services/mizu/components/python"
	"github.com/sdslabs/gasper/services/mizu/components/static"
)

type componentBinding struct {
	validator func(c *gin.Context)
}

var componentMap = map[string]*componentBinding{
	"nodejs": &componentBinding{
		validator: nodejs.Validator,
	},
	"php": &componentBinding{
		validator: php.Validator,
	},
	"python": &componentBinding{
		validator: python.Validator,
	},
	"static": &componentBinding{
		validator: static.Validator,
	},
}
