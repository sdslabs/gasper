package mizu

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/types"
	"github.com/sdslabs/SWS/services/mizu/components/nodejs"
	"github.com/sdslabs/SWS/services/mizu/components/php"
	"github.com/sdslabs/SWS/services/mizu/components/python"
	"github.com/sdslabs/SWS/services/mizu/components/static"
)

type componentBinding struct {
	validator func(c *gin.Context)
	pipeline  func(data map[string]interface{}) types.ResponseError
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
