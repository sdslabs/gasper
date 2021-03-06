package configs

import (
	"fmt"
	"strings"
)

// CreateStaticContainerConfig takes the name of the static application
// and generates the container level config for the same
// Location is the path of index.html or index.htm, leave empty if same
func CreateStaticContainerConfig(name, index string) string {
	path := fmt.Sprintf("%s/%s", GasperConfig.ProjectRoot, name)
	return fmt.Sprintf(`
server {
	listen 80 default_server;
	listen [::]:80 default_server;
	server_name  _;

	sendfile on;
	sendfile_max_chunk 1m;
	tcp_nopush on;

	access_log  /var/log/nginx/access.log  main;
	error_log   /var/log/nginx/error.log   warn;

	location / {
		root   %s/;
		index  %s index.html;
	}

	error_page   500 502 503 504  /50x.html;
	location = /50x.html {
		root   /usr/share/nginx/html;
	}
}
	`, path, index)
}

// CreatePHPContainerConfig takes the name of the PHP application
// and generates the container level config for the same
// Location is the path of index.php, leave empty if same
func CreatePHPContainerConfig(name, index string) string {
	path := fmt.Sprintf("%s/%s", GasperConfig.ProjectRoot, name)

	if strings.Contains(index, "/") {
		subDirs := strings.Split(index, "/")
		index = subDirs[len(subDirs)-1]
		path = fmt.Sprintf("%s/%s", path, strings.Join(subDirs[:len(subDirs)-1], "/"))
	}

	return fmt.Sprintf(`
server {
	listen 80 default_server;
	listen [::]:80 default_server;
	
	server_name _;

	sendfile on;
	sendfile_max_chunk 1m;
	tcp_nopush on;

	access_log  /var/log/nginx/access.log;
	error_log   /var/log/nginx/error.log   warn;

	root %s/;
	index %s index.php index.html;

	location / {
		try_files $uri $uri/ /index.php?q=$uri&$args;
	}

	error_page 500 502 503 504 /50x.html;
	location = /50x.html {
		root /var/lib/nginx/html;
	}

	# pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
	location ~ \.php$ {
		try_files $uri =404;
		fastcgi_split_path_info ^(.+\.php)(/.+)$;
		fastcgi_pass  127.0.0.1:9000;
		fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
		fastcgi_param SCRIPT_NAME $fastcgi_script_name;
		fastcgi_index index.php;
		include fastcgi_params;
	}

	location ~* \.(jpg|jpeg|gif|png|css|js|ico|xml)$ {
		expires 5d;
	}
}
	`, path, index)
}
