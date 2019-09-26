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
	appName := strings.Split(r.Host, ".")[0]
	appURL, err := redis.FetchAppServer(appName)
	if err != nil {
		w.Write([]byte("Could not resolve the requested host."))
		w.WriteHeader(404)
		return
	}

	reverseProxy(w, r, appURL)
	w.WriteHeader(200)
}

// BuildEnraiServer sets up the gorilla multiplexer to handle different subdomains
func BuildEnraiServer(service string) *http.Server {
	enraiConfig := configs.ServiceConfig[service].(map[string]interface{})
	domain := configs.SWSConfig["domain"].(string)

	router := mux.NewRouter().StrictSlash(true)
	host := fmt.Sprintf(`{_:.+}.%s`, domain)
	router.HandleFunc("/", subdomainRootHandler).Host(host)

	server := &http.Server{
		Addr:         enraiConfig["port"].(string),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
