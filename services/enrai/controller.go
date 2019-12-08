package enrai

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

const (
	// DefaultServiceName is the name of the service proxying HTTP connections
	DefaultServiceName = types.Enrai

	// SSLServiceName is the name of the service proxying HTTPS connections
	SSLServiceName = types.EnraiSSL
)

// storage stores the reverse proxy records in the form of Key : Value pairs
// with Application Name as the key and its URL(IP:Port) as the value
var storage = types.NewRecordStorage()

// balancedInstances are the services for which Enrai load balances the request among multiple instances
var balancedInstances = []string{
	types.Kaze,
}

// kazeBalancer load balances requests among multiple kaze instances
var kazeBalancer = types.NewLoadBalancer()

var (
	// Root domain name for validating host names
	rootDomain = fmt.Sprintf(".%s", configs.GasperConfig.Domain)

	// Root domain name with port for validating host names
	rootDomainWithPort = fmt.Sprintf("%s:%d", rootDomain, configs.ServiceConfig.Enrai.Port)
)

// reverseProxy sets up the reverse proxy from the given domain to the target IP
func reverseProxy(c *gin.Context) {
	if strings.HasSuffix(c.Request.Host, rootDomain) || strings.HasSuffix(c.Request.Host, rootDomainWithPort) {
		name := strings.Split(c.Request.Host, ".")[0]
		var target string
		var success bool

		if utils.Contains(balancedInstances, name) {
			target, success = kazeBalancer.Get()
		} else {
			target, success = storage.Get(name)
		}

		if !success {
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
	router.Use(gin.Recovery())
	router.NoRoute(reverseProxy)
	return router
}
