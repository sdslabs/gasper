# Deploying a Node.js Application with MySQL database

This example shows how to deploy a [node.js](https://nodejs.org/en/) application [MySQL](https://www.mysql.com/) database via Gasper

Lets use a [sample application](https://github.com/sdslabs/gasper-sample-nodejs-db) for demonstration which runs on **port 3005** 

!!!warning "Prerequisites"
    * You have [Master](/configurations/master/),  [AppMaker](/configurations/appmaker/) and [DbMaker](/configurations/dbmaker/) up and running
    * You have [DbMaker MySQL Plugin](/configurations/dbmaker/#mysql-configuration) enabled
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token


## Create a MySQL database via Gasper

```bash
$ curl -X POST \
  http://localhost:3000/dbs/mysql \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "nodetestapp",
	"password": "nodetestapp"
}'

{
    "name": "nodetestapp",
    "password": "nodetestapp",
    "user": "nodetestapp",
    "instance_type": "database",
    "language": "mysql",
    "host_ip": "10.43.3.24",
    "port": 33061,
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Now, change the host and port in config file of 

## Deploy using Build and Run Commands

Note: Here host and port should be the host ip and port of the db you just hosted.

```bash
$ curl -X POST \
  http://localhost:3000/apps/nodejs \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"samplenode",
"password":"samplenode",
"git": {
	"repo_url": "https://github.com/sdslabs/gasper-sample-nodejs-db"
},
"context":{
    "index":"index.js",
    "port": 3005,
    "build": ["npm install"],
    "run": ["node index.js"]
},
"env": ["DB_HOST":"10.43.3.24", "DB_PORT":"33061"]
}'

{
    "name": "samplenode",
    "password": "samplenode",
    "git": {
        "repo_url": "https://github.com/sdslabs/gasper-sample-nodejs-db"
    },
    "context": {
        "index": "index.js",
        "port": 3005,
        "rc_file": false,
        "build": [
            "npm install"
        ],
        "run": [
            "node index.js"
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

Have a look at the [run commands file](https://github.com/sdslabs/gasper-sample-nodejs/blob/master/Gasperfile.txt) for the above [sample application](https://github.com/sdslabs/gasper-sample-nodejs)

```bash
$ curl -X POST \
  http://localhost:3000/apps/nodejs \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"samplenode",
"password":"samplenode",
"git": {
	"repo_url": "https://github.com/sdslabs/gasper-sample-nodejs-db"
},
"context":{
    "index":"index.js",
    "port": 3005,
    "rc_file": true
}
}'

{
    "name": "samplenode",
    "password": "samplenode",
    "git": {
        "repo_url": "https://github.com/sdslabs/gasper-sample-nodejs-db"
    },
    "context": {
        "index": "index.js",
        "port": 3005,
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
