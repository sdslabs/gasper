# Deploying an Advanced PHP Application

This example shows how to deploy a PHP application which uses [composer](https://getcomposer.org/) for managing dependencies

Lets use an [advanced PHP application](https://github.com/alphadose/MVC-Project) for demonstration

!!!warning "Prerequisites"
    * You have [Master](/configurations/master/) and [AppMaker](/configurations/appmaker/) up and running
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token


## Deploy using Build and Run Commands

```bash
$ curl -X POST \
  http://localhost:3000/apps/php \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"advancedphp",
"password":"advancedphp",
"git": {
	"repo_url": "https://github.com/alphadose/MVC-Project"
},
"context":{
    "index":"public/index.php",
    "build": ["composer install"],
    "port": 80
}
}'

{
    "name": "advancedphp",
    "password": "advancedphp",
    "git": {
        "repo_url": "https://github.com/alphadose/MVC-Project"
    },
    "context": {
        "index": "public/index.php",
        "port": 80,
        "rc_file": false,
        "build": [
            "composer install"
        ]
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
    "container_id": "f37749b727988833dda70714539ee1ce7f167abe66d78300553f6843a8af39e2",
    "container_port": 50475,
    "language": "php",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 advancedphp@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:50475` 

## Deploy using [Run Commands File](/configurations/global/#run-commands-file)

Have a look at the [run commands file](https://github.com/alphadose/MVC-Project/blob/master/Gasperfile.txt) for the above [sample application](https://github.com/alphadose/MVC-Project)

```bash
$ curl -X POST \
  http://localhost:3000/apps/php \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"advancedphp",
"password":"advancedphp",
"git": {
	"repo_url": "https://github.com/alphadose/MVC-Project"
},
"context":{
    "index":"public/index.php",
    "rc_file": true
}
}'

{
    "name": "advancedphp",
    "password": "advancedphp",
    "git": {
        "repo_url": "https://github.com/alphadose/MVC-Project"
    },
    "context": {
        "index": "public/index.php",
        "port": 80,
        "rc_file": true
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
    "container_id": "d4a54b0800eb8e8bbcea007275746180e5c193b23fc0e1f4f184abf9b984165b",
    "container_port": 51223,
    "language": "php",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 advancedphp@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:51223`
