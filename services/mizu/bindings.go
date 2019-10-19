package mizu

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/services/mizu/components/nodejs"
	"github.com/sdslabs/gasper/services/mizu/components/php"
	"github.com/sdslabs/gasper/services/mizu/components/python"
	"github.com/sdslabs/gasper/services/mizu/components/static"
	"github.com/sdslabs/gasper/types"
)

type componentBinding struct {
	validator func(c *gin.Context)
	pipeline  func(data types.M) types.ResponseError
}

var componentMap = map[string]*componentBinding{
	"nodejs": &componentBinding{
		validator: nodejs.Validator,
		pipeline:  nodejs.Pipeline,
	},
	"php": &componentBinding{
		validator: php.Validator,
		pipeline:  php.Pipeline,
	},
	"python": &componentBinding{
		validator: python.Validator,
		pipeline:  python.Pipeline,
	},
	"static": &componentBinding{
		validator: static.Validator,
		pipeline:  static.Pipeline,
	},
}
