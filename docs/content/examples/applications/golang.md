# Deploying a Golang Application

This example shows how to deploy a [golang](https://golang.org/) application

Lets use a [sample application](https://github.com/sdslabs/gasper-sample-golang) for demonstration which runs on **port 8000** 

!!!warning "Prerequisites"
    * You have [Master](/configurations/master/) and [AppMaker](/configurations/appmaker/) up and running
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token


## Deploy using Build and Run Commands

```bash
$ curl -X POST \
  http://localhost:3000/apps/golang \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"samplego",
"password":"samplego",
"git": {
	"repo_url": "https://github.com/sdslabs/gasper-sample-golang"
},
"context":{
    "index":"main.go",
    "port": 8000,
    "run": ["go run main.go"]
}
}'

{
    "name": "samplego",
    "password": "samplego",
    "git": {
        "repo_url": "https://github.com/sdslabs/gasper-sample-golang"
    },
    "context": {
        "index": "main.go",
        "port": 8000,
        "rc_file": false,
        "run": [
            "go run main.go"
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
    "docker_image": "sdsws/golang:1.1",
    "container_id": "e0d2b784cab9c6cc4c360c81953502d757447a70a7d84bc944cc05819d2ee818",
    "container_port": 55147,
    "language": "golang",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 samplego@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:55147` 

## Deploy using [Run Commands File](/configurations/global/#run-commands-file)

Have a look at the [run commands file](https://github.com/sdslabs/gasper-sample-golang/blob/master/Gasperfile.txt) for the above [sample application](https://github.com/sdslabs/gasper-sample-golang)

```bash
$ curl -X POST \
  http://localhost:3000/apps/golang \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"samplego",
"password":"samplego",
"git": {
	"repo_url": "https://github.com/sdslabs/gasper-sample-golang"
},
"context":{
    "index":"main.go",
    "port": 8000,
    "rc_file": true
}
}'

{
    "name": "samplego",
    "password": "samplego",
    "git": {
        "repo_url": "https://github.com/sdslabs/gasper-sample-golang"
    },
    "context": {
        "index": "main.go",
        "port": 8000,
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
    "docker_image": "sdsws/golang:1.1",
    "container_id": "0c4b1ec05fe65fcb0b3ef168244d38ae9fab4d0bc22e2e0d5a39badeedce31e7",
    "container_port": 55229,
    "language": "golang",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 samplego@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:55229`
