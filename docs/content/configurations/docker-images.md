# Docker Images Configuration

The docker images used by Gasper for creating application containers and databases are defined in this section

```toml
###################################
#   Docker Images Configuration   #
###################################

[images]
static = "docker.io/sdslabs/static:latest"
php = "docker.io/sdslabs/php:latest"
nodejs = "docker.io/sdslabs/node:latest"
python2 =  "docker.io/sdslabs/python2:latest"
python3 = "docker.io/sdslabs/python3:latest"
golang = "docker.io/sdslabs/golang:latest"
ruby = "docker.io/sdslabs/ruby:latest"
rust = "docker.io/sdslabs/rust:latest"
mysql = "docker.io/mysql:latest"
mongodb = "docker.io/sdslabs/alpine-mongo:latest"
postgresql = "docker.io/postgres:latest"
redis = "docker.io/redis:6.0-rc3-alpine3.11"
```

You can replace the above default images and plug in your own docker images but make sure that each image has a **blocking CMD call** at the end of its corresponding dockerfile such as **CMD tail -f /dev/null**

For reference, you can check out the [dockerfiles](https://github.com/sdslabs/gasper-dockerfiles) for the default images used by Gasper
