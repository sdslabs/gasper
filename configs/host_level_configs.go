package configs

import (
	"fmt"
)

// CreateApacheHostConfig takes the port (string) and assigns returns a conf file for host
func CreateApacheHostConfig(port string) string {
	return fmt.Sprintf(`
<VirtualHost *:80>
    # Set the below ServerAlias
    # eg: *.mysite.com for the static container
    ServerAlias *.%s
    # Reverse proxy for the port pointing to container
    ProxyPass / http://localhost:%s/
    ProxyPassReverse / http://localhost:%s/
    # To set the HOSTNAME received by the container as the ServerName, not 'localhost:port'
    RequestHeader set Host %%{HTTP_HOST}
    ProxyPreserveHost On
    # Error log
    ErrorLog /var/log/apache2/static.error.log
    LogLevel warn
    CustomLog /var/log/apache2/static.access.log combined
</VirtualHost>
    `, SWSConfig["domain"].(string), port, port)
}
