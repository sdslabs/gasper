package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SDS/utils"
)

// createStaticPage handles making of new static pages
func createStaticPage(json staticAppConfig) (gin.H, utils.Error) {
	return gin.H{}, utils.Error{
		Code: 200,
		Err:  nil,
	}
}
