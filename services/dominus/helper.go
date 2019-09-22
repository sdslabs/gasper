package dominus

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/configs"
)

func trimURLPath(c *gin.Context) {
	urlPathSlice := strings.Split(c.Request.URL.Path, "/")
	if len(urlPathSlice) >= 2 {
		c.Request.URL.Path = fmt.Sprintf("/%s", strings.Join(urlPathSlice[2:], "/"))
		c.Next()
	} else {
		c.JSON(404, gin.H{
			"message": "Page not found",
		})
	}
}

func reverseProxy(c *gin.Context, target string) {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = target
		req.Host = target
		req.Header["dominus-secret"] = []string{configs.SWSConfig["secret"].(string)}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}
