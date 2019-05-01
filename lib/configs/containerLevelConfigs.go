package configs

import (
	"fmt"

	"github.com/sdslabs/SWS/lib/utils"
)

// CreateStaticContainerConfig takes the name of the static application
// and generates the container level config for the same
// Location is the path of index.html or index.htm, leave empty if same
func CreateStaticContainerConfig(name string, appContext map[string]interface{}) string {
	path := fmt.Sprintf("%s/%s", name, appContext["index"].(string))
	return fmt.Sprintf(`
server {
	listen       80;
	server_name  %s.%s;

	access_log  /var/log/nginx/%s.access.log  main;
	error_log   /var/log/nginx/%s.error.log   warn;

	location / {
		root   /SWS/%s;
		index  index.html index.htm;
	}

	error_page   500 502 503 504  /50x.html;
	location = /50x.html {
		root   /usr/share/nginx/html;
	}
}
	`, name, utils.SWSConfig["domain"].(string), name, name, path)
}

// CreatePHPContainerConfig takes the name of the PHP application
// and generates the container level config for the same
// Location is the path of index.php, leave empty if same
func CreatePHPContainerConfig(name string, appContext map[string]interface{}) string {
	path := fmt.Sprintf("%s/%s", name, appContext["index"].(string))
	return fmt.Sprintf(`
server {
	listen 80;
	listen [::]:80;
	server_name %s.%s;

	access_log  /var/log/nginx/%s.access.log  main;
	error_log   /var/log/nginx/%s.error.log   warn;

	root /SWS/%s;
	index index.php;

	location / {
		try_files  / /index.php?;
	}

	location ~ \.php$ {
		try_files  =404;
		fastcgi_split_path_info ^(.+\.php)(/.+)$;
		fastcgi_pass unix:/var/run/php/php7.0-fpm.sock;
		fastcgi_param SCRIPT_FILENAME ;
		fastcgi_index index.php;
		include fastcgi_params;
	}

	error_page   500 502 503 504  /50x.html;
	location = /50x.html {
		root   /usr/share/nginx/html;
	}
}
`, name, utils.SWSConfig["domain"].(string), name, name, path)
}

// CreateNodeContainerConfig takes the name of the node app
// and port and generated the config for the same
func CreateNodeContainerConfig(name string, appContext map[string]interface{}) string {
	return fmt.Sprintf(`
server {
    listen 80;
    server_name %s.%s;

    location / {
    	proxy_set_header   X-Forwarded-For $remote_addr;
    	proxy_set_header   Host $http_host;
    	proxy_pass         http://127.0.0.1:%s;
	}
}
`, name, utils.SWSConfig["domain"].(string), appContext["port"].(string))
}

// CreatePythonContainerConfig takes the name of the Python app
// and port and generated the config for the same
func CreatePythonContainerConfig(name string, appContext map[string]interface{}) string {
	return fmt.Sprintf(`
server {
    listen 80;
    server_name %s.%s;

    location / {
    	proxy_set_header   X-Forwarded-For $remote_addr;
    	proxy_set_header   Host $http_host;
    	proxy_pass         http://127.0.0.1:%s;
	}
}
`, name, utils.SWSConfig["domain"].(string), appContext["port"].(string))
}
