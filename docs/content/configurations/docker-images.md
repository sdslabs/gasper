# Docker Images Configuration

The docker images used by Gasper for creating application containers and databases are defined in this section

```toml
###################################
#   Docker Images Configuration   #
###################################

[images]
static = "docker.io/sdsws/static:2.0"
php = "docker.io/sdsws/php:3.0"
nodejs = "docker.io/sdsws/node:2.1"
python2 =  "docker.io/sdsws/python2:1.1"
python3 = "docker.io/sdsws/python3:1.1"
golang = "docker.io/sdsws/golang:1.1"
ruby = "docker.io/sdsws/ruby:1.0"
mysql = "docker.io/wangxia/alpine-mysql"
mongodb = "docker.io/sdsws/alpine-mongo"
postgresql = "docker.io/postgres:12.2-alpine"
redis = "docker.io/redis:6.0-rc3-alpine3.11"
```

You can replace the above default images and plug in your own docker images but make sure that each image has a **blocking CMD call** at the end of its corresponding dockerfile such as **CMD tail -f /dev/null**

For reference, you can check out the [dockerfiles](https://github.com/sdslabs/gasper-dockerfiles) for the default images used by Gasper
