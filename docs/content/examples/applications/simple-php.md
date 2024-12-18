# Deploying a Simple PHP Application

This example shows how to deploy a simple PHP application

Lets use a [sample PHP application](https://github.com/sdslabs/gasper-sample-php) for demonstration

!!!warning "Prerequisites"
    * You have [Master](/configurations/master/) and [AppMaker](/configurations/appmaker/) up and running
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token

```bash
$ curl -X POST \
  http://localhost:3000/apps/php \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"simplephp",
"password":"simplephp",
"git": {
	"repo_url": "https://github.com/sdslabs/gasper-sample-php"
},
"context":{
    "index":"index.php",
    "port":80
}
}'

{
    "name": "simplephp",
    "password": "simplephp",
    "git": {
        "repo_url": "https://github.com/sdslabs/gasper-sample-php"
    },
    "context": {
        "index": "index.php",
        "port": 80,
        "rc_file": false
    },
    "resources": {
        "memory": 0.5,
        "cpu": 0.25
    },
    "name_servers": [
        "192.168.108.121",
        "192.168.108.122",
        "10.43.3.24"
    ],
    "docker_image": "sdsws/php:3.0",
    "container_id": "c447c03399e5b23b860c6bfd932fa6a7f93e9ff6d7001d0cd4064f1554752cc3",
    "container_port": 49599,
    "language": "php",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 simplephp@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:49599` 

