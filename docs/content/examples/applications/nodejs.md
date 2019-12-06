# Deploying a Node.js Application

This example shows how to deploy a [node.js](https://nodejs.org/en/) application

Lets use a [sample application](https://github.com/sdslabs/sample-nodejs) for demonstration which runs on **port 3000** 

!!!warning "Prerequisites"
    * You have [Kaze](/configurations/kaze/) and [Mizu](/configurations/mizu/) up and running
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token


## Deploy using Build and Run Commands

```bash
$ curl -X POST \
  http://localhost:3000/apps/nodejs \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"samplenode",
"password":"samplenode",
"git": {
	"repo_url": "https://github.com/sdslabs/sample-nodejs"
},
"context":{
    "index":"main.js",
    "port": 3000,
    "build": ["npm install"],
    "run": ["node main.js"]
}
}'

{
    "name": "samplenode",
    "password": "samplenode",
    "git": {
        "repo_url": "https://github.com/sdslabs/sample-nodejs"
    },
    "context": {
        "index": "main.js",
        "port": 3000,
        "rc_file": false,
        "build": [
            "npm install"
        ],
        "run": [
            "node main.js"
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
    "docker_image": "sdsws/node:2.1",
    "container_id": "dc04aea7dbef287b5bfa597120773c4ff5b5309d3a39235055ff80e9ffbee00f",
    "container_port": 51952,
    "language": "nodejs",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 samplenode@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:51952` 

## Deploy using [Run Commands File](/configurations/global/#run-commands-file)

Have a look at the [run commands file](https://github.com/sdslabs/sample-nodejs/blob/master/Gasperfile.txt) for the above [sample application](https://github.com/sdslabs/sample-nodejs)

```bash
$ curl -X POST \
  http://localhost:3000/apps/nodejs \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"samplenode",
"password":"samplenode",
"git": {
	"repo_url": "https://github.com/sdslabs/sample-nodejs"
},
"context":{
    "index":"main.js",
    "port": 3000,
    "rc_file": true
}
}'

{
    "name": "samplenode",
    "password": "samplenode",
    "git": {
        "repo_url": "https://github.com/sdslabs/sample-nodejs"
    },
    "context": {
        "index": "main.js",
        "port": 3000,
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
    "docker_image": "sdsws/node:2.1",
    "container_id": "5e025f17c9d5c11f93609f7d019b3efbdc44ccef598a6e6564973da895e5e366",
    "container_port": 51720,
    "language": "nodejs",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 samplenode1@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:51720`
