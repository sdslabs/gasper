package genproxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
)

const (
	// DefaultServiceName is the name of the service proxying HTTP connections
	DefaultServiceName = types.GenProxy

	// SSLServiceName is the name of the service proxying HTTPS connections
	SSLServiceName = types.GenProxySSL
)

var (
	// storage stores the reverse proxy records in the form of Key : Value pairs
	// with Application Name as the key and its URL(IP:Port) as the value
	storage = types.NewProxyStorage()

	// balancedInstances are the services for which GenProxy load balances the
	// request among multiple instances
	balancedInstances = []string{
		types.Master,
		"gasper",
	}

	// masterBalancer load balances requests among multiple master instances
	masterBalancer = types.NewLoadBalancer()

	// Root domain name for validating host names
	rootDomain = fmt.Sprintf(".%s", configs.GasperConfig.Domain)

	// Root domain name with port for validating host names
	rootDomainWithPort = fmt.Sprintf("%s:%d", rootDomain, configs.ServiceConfig.GenProxy.Port)
)

// reverseProxy sets up the reverse proxy from the given domain to the target IP
func reverseProxy(c *gin.Context) {
	if !strings.HasSuffix(c.Request.Host, rootDomain) && !strings.HasSuffix(c.Request.Host, rootDomainWithPort) {
		c.AbortWithStatusJSON(403, gin.H{
			"success": false,
			"message": "Incorrect root domain",
		})
		return
	}

	name := strings.Split(c.Request.Host, ".")[0]
	var proxy *types.ProxyInfo
	var success bool

	if utils.Contains(balancedInstances, name) {
		proxy, success = masterBalancer.Get()
	} else {
		proxy, success = storage.Get(name)
	}

	if !success {
		c.AbortWithStatusJSON(503, gin.H{
			"success": false,
			"message": "No such application exists",
		})
		return
	}
	proxy.Serve(c)
}

// NewService returns a new instance of the current microservice
func NewService() http.Handler {
	// router is the main routes handler for the current microservice package
	router := gin.New()
	router.Use(gin.Recovery())
	router.NoRoute(reverseProxy)
	return router
}
