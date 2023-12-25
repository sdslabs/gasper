# Deploying a Static Website

This example shows how to deploy a static website

Lets use the [hangman game](https://github.com/sdslabs/hangman-js-game) for demonstration

!!!warning "Prerequisites"
    * You have [Master](/configurations/master/) and [AppMaker](/configurations/appmaker/) up and running
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token

```bash
$ curl -X POST \
  http://localhost:3000/apps/static \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"static",
"password":"static",
"git": {
	"repo_url": "https://github.com/sdslabs/hangman-js-game"
},
"context":{
    "index":"hangman.html",
    "port":80
}
}'

{
    "name": "static",
    "password": "static",
    "git": {
        "repo_url": "https://github.com/sdslabs/hangman-js-game"
    },
    "context": {
        "index": "hangman.html",
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
    "docker_image": "sdsws/static:2.0",
    "container_id": "a05900527ad4b7175be438d8d28707cda39df3b94806d35f92949fd0b3d134db",
    "container_port": 65499,
    "language": "static",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 static@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:65499` 
