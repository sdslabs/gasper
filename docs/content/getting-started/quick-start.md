# Quick Start

## Grab the latest binary
Assuming you have the [dependencies](/getting-started/dependencies) installed, head over to Gasper's [releases](https://github.com/sdslabs/gasper/releases) page and grab the latest binary according to your operating system and system architecture.

## Extract the downloaded content
After downloading, unzip the tar file

```bash
$ tar -xf gasper_platform_arch.tar.gz
```

After extraction, the extracted directory should have the `gasper binary` and `config.toml`, the configuration file

```bash
$ cd gasper_platform_arch
$ ls
gasper
config.toml
```

## Run Gasper
Run Gasper by executing the binary with the provided configuration
```bash
$ ./gasper --conf ./config.toml
```

!!!warning
    Make sure that Docker, Redis and MongoDB are running on your system before executing the above command

## Login and Token Retrieval
After Gasper is up and successfully running, lets deploy a sample application using [curl](https://curl.haxx.se/)

To do that first we need to login and obtain a [JWT](https://jwt.io/) (JSON Web Token)

```bash
$ curl -X POST \
  http://localhost:3000/auth/login \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "anish.mukherjee1996@gmail.com",
    "password": "alphadose"
  }'

{
    "code": 200,
    "expire": "2019-12-04T22:05:41+05:30",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJhbHBoYWRvc2VAZ21haWwuY29tIiwiZXhwIjoxNTc1NDc3MzQxLCJvcmlnX2lhdCI6MTU3NTQ3Mzc0MSwidXNlcm5hbWUiOiJhbHBoYWRvc2UifQ.Io0txryVH8zR6JfZ0iey86474oZl8gNwo4HjKgZl2s8"
}
```

!!!note
    If you have made any changes in the [admin section](https://github.com/sdslabs/gasper/blob/develop/config.sample.toml#L36) of `config.toml` then change the payload (email and password) of the above request accordingly

## Application Deployment
The **token** obtained from the above JSON response is our required JWT

We will now use that **token** in the **Authorization Header** to deploy a [Sample PHP application](https://github.com/sdslabs/sample-php)
The format for using the token in the request header is `Authorization: Bearer {{token}}`

```bash
$ curl -X POST \
  http://localhost:3000/apps/php \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJhbHBoYWRvc2VAZ21haWwuY29tIiwiZXhwIjoxNTc1NDc4MTc5LCJvcmlnX2lhdCI6MTU3NTQ3NDU3OSwidXNlcm5hbWUiOiJhbHBoYWRvc2UifQ.XKxKmC5mrSwHq3RGmTGqiAcQreVQjd9S-DMxw8ZN1k0' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"test",
"password":"test",
"git": {
	"repo_url": "https://github.com/sdslabs/sample-php"
},
"context":{
    "index":"index.php"
}
}'

{
    "name": "test",
    "password": "test",
    "git": {
        "repo_url": "https://github.com/sdslabs/sample-php"
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
    "container_id": "fe04f8d7cbbdfa100ac9f03c8bdcec7b3d3246aa189dc0264c7d2af1cb92308b",
    "container_port": 64128,
    "language": "php",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 test@10.43.3.24",
    "owner": "alphadose@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:64128`

You should get the message `Hello World` in your browser marking the end of this tutorial

!!!question "Where to go next?"
    You can either have a look at more [examples](/examples/login) or how to [configure and setup](/configurations/overview) Gasper to your liking
