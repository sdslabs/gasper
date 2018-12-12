package configs

import (
	"fmt"

	"github.com/sdslabs/SWS/lib/utils"
)

// CreateStaticContainerConfig takes the name of the static application
// and generates the container level config for the same
func CreateStaticContainerConfig(name string) string {
	return fmt.Sprintf(`
server {
	listen       80;
	server_name  %s.%s;

	access_log  /var/log/nginx/app.access.log  main;
	error_log   /var/log/nginx/app.error.log   warn;

	location / {
		root   /SWS/app/;
		index  index.html index.htm;
	}

	error_page   500 502 503 504  /50x.html;
	location = /50x.html {
		root   /usr/share/nginx/html;
	}
}
	`, name, utils.SWSConfig.Domain)
}
