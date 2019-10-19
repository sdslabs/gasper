package enrai

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/redis"
)

// ServiceName is the name of the current microservice
const ServiceName = "enrai"

// reverseProxy sets up the reverse proxy from the given domain
// to the target IP
func reverseProxy(c *gin.Context) {
	rootDomain := fmt.Sprintf(".%s", configs.GasperConfig.Domain)
	rootDomainWithPort := fmt.Sprintf("%s:%d", rootDomain, configs.ServiceConfig.Enrai.Port)

	if strings.HasSuffix(c.Request.Host, rootDomain) || strings.HasSuffix(c.Request.Host, rootDomainWithPort) {
		target, err := redis.FetchAppServer(strings.Split(c.Request.Host, ".")[0])

		if err != nil {
			c.AbortWithStatusJSON(503, gin.H{
				"success": false,
				"message": "No such application exists",
			})
			return
		}

		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
			req.Host = target
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
		return
	}

	c.AbortWithStatusJSON(403, gin.H{
		"success": false,
		"message": "Incorrect root domain",
	})
}

// NewService returns a new instance of the current microservice
func NewService() http.Handler {
	// router is the main routes handler for the current microservice package
	router := gin.New()
	router.NoRoute(reverseProxy)
	return router
}
