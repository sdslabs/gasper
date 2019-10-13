package dominus

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/gasper/configs"
)

func trimURLPath(length int) gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPathSlice := strings.Split(c.Request.URL.Path, "/")
		if len(urlPathSlice) >= length {
			c.Request.URL.Path = fmt.Sprintf("/%s", strings.Join(urlPathSlice[length:], "/"))
			c.Next()
		} else {
			c.AbortWithStatusJSON(404, gin.H{
				"message": "Page not found",
				"success": false,
			})
		}
	}
}

func reverseProxy(c *gin.Context, target string) {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = target
		req.Host = target
		req.Header["dominus-secret"] = []string{configs.GasperConfig.Secret}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}
