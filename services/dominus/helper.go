package dominus

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/SWS/lib/configs"
)

func reverseProxy(c *gin.Context, target string) {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = target
		req.Host = target
		req.Header["dominus-secret"] = []string{configs.SWSConfig["secret"].(string)}
		if req.Method == "POST" {
			req.URL.Path = ""
		}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}
