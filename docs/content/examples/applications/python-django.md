# Deploying a Python Django Application

This example shows how to deploy a [python django](https://www.djangoproject.com/) application

Lets use a [sample application](https://github.com/sdslabs/sample-django) for demonstration which runs on **port 8000** 

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
"name":"sampledjango",
"password":"sampledjango",
"git": {
	"repo_url": "https://github.com/sdslabs/sample-django"
},
"context":{
    "index":"todo/manage.py",
    "port": 8000,
    "build": ["pip install -r requirements.txt", "python todo/manage.py migrate"],
    "run": ["python todo/manage.py runserver 0.0.0.0:8000"]
}
}'

{
    "name": "sampledjango",
    "password": "sampledjango",
    "git": {
        "repo_url": "https://github.com/sdslabs/sample-django"
    },
    "context": {
        "index": "todo/manage.py",
        "port": 8000,
        "rc_file": false,
        "build": [
            "pip install -r requirements.txt",
            "python todo/manage.py migrate"
        ],
        "run": [
            "python todo/manage.py runserver 0.0.0.0:8000"
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
    "container_id": "8f2a04bb54f0b90a911b05d3fb1ae73ff240c2e5a5093609d393f7c426de4755",
    "container_port": 53358,
    "language": "python3",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 sampledjango@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:53358` 

## Deploy using [Run Commands File](/configurations/global/#run-commands-file)

Have a look at the [run commands file](https://github.com/sdslabs/sample-django/blob/master/Gasperfile.txt) for the above [sample application](https://github.com/sdslabs/sample-django)

```bash
$ curl -X POST \
  http://localhost:3000/apps/python3 \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"sampledjango2",
"password":"sampledjango2",
"git": {
	"repo_url": "https://github.com/sdslabs/sample-django"
},
"context":{
    "index":"todo/manage.py",
    "port": 8000,
    "rc_file": true
}
}'

{
    "name": "sampledjango",
    "password": "sampledjango",
    "git": {
        "repo_url": "https://github.com/sdslabs/sample-django"
    },
    "context": {
        "index": "todo/manage.py",
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
    "docker_image": "sdsws/python3:1.1",
    "container_id": "48ada540a6296b184470eb192e4b543f195ab5f67615ec12318d3e8d01e05edf",
    "container_port": 53672,
    "language": "python3",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 sampledjango@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:53672`

!!!info
    HTTP request to the above URL endpoint `localhost:3000/apps/python3` runs the application in a **Python 3** environment. If you want to run your application in a **Python 2** environment then change the URL endpoint to `localhost:3000/apps/python2`
