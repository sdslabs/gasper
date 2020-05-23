package types

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

// ProxyInfo is a container for establishing a reverse-proxy connection
type ProxyInfo struct {
	host       string
	connection *httputil.ReverseProxy
}

// Serve establishes a reverse proxy connection
func (proxy *ProxyInfo) Serve(c *gin.Context) {
	proxy.connection.ServeHTTP(c.Writer, c.Request)
}

// UpdateDirector updates the endpoint in case of any change in the system
func (proxy *ProxyInfo) UpdateDirector(host string) {
	proxy.host = host
	proxy.connection.Director = func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = host
		req.Host = host
	}
}

// NewProxyInfo returns a new ProxyInfo container
func NewProxyInfo(host string) *ProxyInfo {
	return &ProxyInfo{
		host: host,
		connection: &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = host
				req.Host = host
			},
		},
	}
}
