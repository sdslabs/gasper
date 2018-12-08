package configs

import (
	"fmt"

	"github.com/sdslabs/SWS/lib/utils/json"
)

// CreateStaticContainerConfig takes the name of the static application
// and generates the container level config for the same
func CreateStaticContainerConfig(name string) string {
	return fmt.Sprintf(`
	server {
		listen       80;
		server_name  %s.%s;

		access_log  /var/log/nginx/%s.access.log  main;
		error_log   /var/log/nginx/%s.error.log   warn;

		location / {
			root   /SWS/%s/;
			index  index.html index.htm;
		}

		error_page   500 502 503 504  /50x.html;
		location = /50x.html {
			root   /usr/share/nginx/html;
		}
	}
	`, name, json.Domain, name, name, name)
}
