# Deploying a Python Flask Application

This example shows how to deploy a [python flask](https://www.palletsprojects.com/p/flask/) application

Lets use a [sample application](https://github.com/sdslabs/sample-flask) for demonstration which runs on **port 5000** 

!!!warning "Prerequisites"
    * You have [Kaze](/configurations/kaze/) and [Mizu](/configurations/mizu/) up and running
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token


## Deploy using Build and Run Commands

```bash
$ curl -X POST \
  http://localhost:3000/apps/python3 \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"sampleflask",
"password":"sampleflask",
"git": {
	"repo_url": "https://github.com/sdslabs/sample-flask"
},
"context":{
    "index":"run.py",
    "port": 5000,
    "build": ["pip install -r requirements.txt"],
    "run": ["flask run --host=0.0.0.0 --port=5000"]
}
}'

{
    "name": "sampleflask",
    "password": "sampleflask",
    "git": {
        "repo_url": "https://github.com/sdslabs/sample-flask"
    },
    "context": {
        "index": "run.py",
        "port": 5000,
        "rc_file": false,
        "build": [
            "pip install -r requirements.txt"
        ],
        "run": [
            "flask run --host=0.0.0.0 --port=5000"
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
    "docker_image": "sdsws/python3:1.1",
    "container_id": "b9521abaa377f5cdf525eb3e3fbe083719f8bee7f8500863b079310f69f4a413",
    "container_port": 52687,
    "language": "python3",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 sampleflask@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:52687` 

## Deploy using [Run Commands File](/configurations/global/#run-commands-file)

Have a look at the [run commands file](https://github.com/sdslabs/sample-flask/blob/master/Gasperfile.txt) for the above [sample application](https://github.com/sdslabs/sample-flask)

```bash
$ curl -X POST \
  http://localhost:3000/apps/python3 \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"sampleflask",
"password":"sampleflask",
"git": {
	"repo_url": "https://github.com/sdslabs/sample-flask"
},
"context":{
    "index":"run.py",
    "port": 5000,
    "rc_file": true
}
}'

{
    "name": "sampleflask",
    "password": "sampleflask",
    "git": {
        "repo_url": "https://github.com/sdslabs/sample-flask"
    },
    "context": {
        "index": "run.py",
        "port": 5000,
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
    "docker_image": "sdsws/python3:1.1",
    "container_id": "574c8b5d8c9e8a14baa10f207723c2083ff28d008b9302a6bb3a6662cb7b06a8",
    "container_port": 52811,
    "language": "python3",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 sampleflask@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:52811`

!!!info
    HTTP request to the above URL endpoint `localhost:3000/apps/python3` runs the application in a **Python 3** environment. If you want to run your application in a **Python 2** environment then change the URL endpoint to `localhost:3000/apps/python2`
