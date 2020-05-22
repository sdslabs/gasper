# Deploying a Rust Application

This example shows how to deploy a [rust](https://www.rust-lang.org/) application

Lets use a [sample application](https://github.com/sdslabs/gasper-sample-rust) for demonstration which runs on **port 3000** 

!!!warning "Prerequisites"
    * You have [Master](/configurations/master/) and [AppMaker](/configurations/appmaker/) up and running
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token


## Deploy using Build and Run Commands

```bash
$ curl -X POST \
  http://localhost:3000/apps/rust \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"samplerust",
"password":"samplerust",
"git": {
    "repo_url":"https://github.com/sdslabs/gasper-sample-rust"
},
"context":{
    "index": "src/main.rs",
    "port": 3000,
    "build" : ["cargo build --release"],
    "run": ["./target/release/gasper-sample-rust"]
}
},
"resources": {
	"memory": 4,
	"cpu": 4
}
}'

{
    "name": "samplerust",
    "password": "samplerust",
    "git": {
        "repo_url": "https://github.com/sdslabs/gasper-sample-rust"
    },
    "context": {
        "index": "src/main.rs",
        "port": 3000,
        "rc_file": false,
        "build": [
            "cargo build --release"
        ],
        "run": [
            "./target/release/gasper-sample-rust"
        ]
    },
    "resources": {
        "memory": 4,
        "cpu": 4
    },
    "name_servers": [
        "8.8.8.8",
        "8.8.4.4"
    ],
    "docker_image": "sdsws/rust:1.0",
    "container_id": "2b9b1f772259c4c6e81aebc3d0e5aca941695bb923ef55652eb73a6b45765c61",
    "container_port": 53341,
    "language": "rust",
    "instance_type": "application",
    "app_url": "samplerust.app.sdslabs.co",
    "host_ip": "192.168.29.250",
    "ssh_cmd": "ssh -p 2222 samplerust@192.168.29.250",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `192.168.29.250:53341`

!!!warning
    The above [sample application](https://github.com/sdslabs/gasper-sample-rust) takes around 3 minutes to build and start hence you need to wait for that duration before hitting the URL in your browser

## Deploy using [Run Commands File](/configurations/global/#run-commands-file)

Have a look at the [run commands file](https://github.com/sdslabs/gasper-sample-rust/blob/master/Gasperfile.txt) for the above [sample application](https://github.com/sdslabs/gasper-sample-rust)

```bash
$ curl -X POST \
  http://localhost:3000/apps/rust \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"samplerust",
"password":"samplerust",
"git": {
    "repo_url":"https://github.com/sdslabs/gasper-sample-rust"
},
"context":{
    "index": "src/main.rs",
    "port": 3000,
    "rc_file": true
}
},
"resources": {
	"memory": 4,
	"cpu": 4
}
}'

{
    "name": "samplerust",
    "password": "samplerust",
    "git": {
        "repo_url": "https://github.com/sdslabs/gasper-sample-rust"
    },
    "context": {
        "index": "src/main.rs",
        "port": 3000,
        "rc_file": true
    },
    "resources": {
        "memory": 4,
        "cpu": 4
    },
    "name_servers": [
        "8.8.8.8",
        "8.8.4.4"
    ],
    "docker_image": "sdsws/rust:1.0",
    "container_id": "917423498cc1d1d344a00069da6d453fdcdd7848d502a43063f852a4bd8afb94",
    "container_port": 53473,
    "language": "rust",
    "instance_type": "application",
    "app_url": "samplerust.app.sdslabs.co",
    "host_ip": "192.168.29.250",
    "ssh_cmd": "ssh -p 2222 samplerust@192.168.29.250",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `192.168.29.250:53473`

!!!warning
    The above [sample application](https://github.com/sdslabs/gasper-sample-rust) takes around 3 minutes to build and start hence you need to wait for that duration before hitting the URL in your browser
