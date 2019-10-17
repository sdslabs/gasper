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

// reverseProxy sets up the reverse proxy from the given
// to the target ip
func reverseProxy(c *gin.Context) {
	hostNameCheck := fmt.Sprintf(".%s", configs.GasperConfig.Domain)
	hostNameWithPortCheck := fmt.Sprintf("%s:%d", hostNameCheck, configs.ServiceConfig.Enrai.Port)

	if strings.HasSuffix(c.Request.Host, hostNameCheck) || strings.HasSuffix(c.Request.Host, hostNameWithPortCheck) {
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
			req.Header["dominus-secret"] = []string{configs.GasperConfig.Secret}
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)

		return
	}

	c.AbortWithStatusJSON(403, gin.H{
		"success": false,
		"message": "Incorrect root domain",
	})
	return
}

// BuildEnraiServer sets up the gorilla multiplexer to handle different subdomains
func BuildEnraiServer() *gin.Engine {
	router := gin.New()
	router.NoRoute(reverseProxy)
	return router
}
