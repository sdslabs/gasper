package enrai

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/redis"
)

// reverseProxy sets up the reverse proxy from the given
// to the target ip
func reverseProxy(w http.ResponseWriter, r *http.Request, target string) {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = target
		req.Host = target
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

// subdomainRootHandler handles the root route of the provided host
// and extracts the url to perform the reverse proxy
func subdomainRootHandler(w http.ResponseWriter, r *http.Request) {
	appURL, err := redis.FetchAppServer(strings.Split(r.Host, ".")[0])
	if err != nil {
		w.WriteHeader(404)
		return
	}
	reverseProxy(w, r, appURL)
}

// BuildEnraiServer sets up the gorilla multiplexer to handle different subdomains
func BuildEnraiServer(service string) *http.Server {
	enraiConfig := configs.ServiceConfig[service].(map[string]interface{})

	router := mux.NewRouter()
	host := fmt.Sprintf(`{_:.+}.%s`, configs.SWSConfig["domain"].(string))
	router.PathPrefix("/").HandlerFunc(subdomainRootHandler).Host(host)

	server := &http.Server{
		Addr:         enraiConfig["port"].(string),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
