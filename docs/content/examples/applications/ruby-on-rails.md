# Deploying a Ruby on Rails Application

This example shows how to deploy a [ruby on rails](https://rubyonrails.org/) application

Lets use a [sample application](https://github.com/sdslabs/ruby-on-rails-sample-app) for demonstration which runs on **port 3000** 

!!!warning "Prerequisites"
    * You have [Kaze](/configurations/kaze/) and [Mizu](/configurations/mizu/) up and running
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token


## Deploy using Build and Run Commands

```bash
$ curl -X POST \
  http://localhost:3000/apps/ruby \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"sampleruby",
"password":"sampleruby",
"git": {
	"repo_url": "https://github.com/sdslabs/ruby-on-rails-sample-app"
},
"context":{
    "index":"bin/rails",
    "port": 3000,
    "build": ["bundle install --without production", "rails db:migrate"],
    "run": ["rails server"] 
},
"resources": {
	"memory": 4,
	"cpu": 4
}
}'

{
    "name": "sampleruby",
    "password": "sampleruby",
    "git": {
        "repo_url": "https://github.com/sdslabs/ruby-on-rails-sample-app"
    },
    "context": {
        "index": "bin/rails",
        "port": 3000,
        "rc_file": false,
        "build": [
            "bundle install --without production",
            "rails db:migrate"
        ],
        "run": [
            "rails server"
        ]
    },
    "resources": {
        "memory": 4,
        "cpu": 4
    },
    "name_servers": [
        "192.168.108.121",
        "192.168.108.122",
        "10.43.3.24"
    ],
    "docker_image": "sdsws/ruby:1.0",
    "container_id": "dd4d4199b81120abe58fb80dca355eba639e1caf8fb37ade02c9a53ee40634a0",
    "container_port": 55673,
    "language": "ruby",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 sampleruby@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:55673`

!!!warning
    The above [sample application](https://github.com/sdslabs/ruby-on-rails-sample-app) takes around 6 minutes to start hence you need to wait for that duration before hitting the URL in your browser

## Deploy using [Run Commands File](/configurations/global/#run-commands-file)

Have a look at the [run commands file](https://github.com/sdslabs/ruby-on-rails-sample-app/blob/master/Gasperfile.txt) for the above [sample application](https://github.com/sdslabs/ruby-on-rails-sample-app)

```bash
$ curl -X POST \
  http://localhost:3000/apps/ruby \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
"name":"sampleruby",
"password":"sampleruby",
"git": {
	"repo_url": "https://github.com/sdslabs/ruby-on-rails-sample-app"
},
"context":{
    "index":"bin/rails",
    "port": 3000,
    "rc_file": true
},
"resources": {
	"memory": 4,
	"cpu": 4
}
}'

{
    "name": "sampleruby",
    "password": "sampleruby",
    "git": {
        "repo_url": "https://github.com/sdslabs/ruby-on-rails-sample-app"
    },
    "context": {
        "index": "bin/rails",
        "port": 3000,
        "rc_file": true
    },
    "resources": {
        "memory": 4,
        "cpu": 4
    },
    "name_servers": [
        "192.168.108.121",
        "192.168.108.122",
        "10.43.3.24"
    ],
    "docker_image": "sdsws/ruby:1.0",
    "container_id": "2e1b2165f93836d8021465802857692f37b51515361e78d13d201fde645d753f",
    "container_port": 56041,
    "language": "ruby",
    "instance_type": "application",
    "host_ip": "10.43.3.24",
    "ssh_cmd": "ssh -p 2222 sampleruby@10.43.3.24",
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```

Note the **host_ip** and **container_port** fields in the above JSON response

You can now access the deployed application by hitting the URL **host_ip:container_port** from your browser

For the above case it will be `10.43.3.24:56041`

!!!warning
    The above [sample application](https://github.com/sdslabs/ruby-on-rails-sample-app) takes around 6 minutes to start hence you need to wait for that duration before hitting the URL in your browser
