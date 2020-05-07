# Docker Images Configuration

The docker images used by Gasper for creating application containers and databases are defined in this section

```toml
###################################
#   Docker Images Configuration   #
###################################

[images]
static = "sdsws/static:2.0"
php = "sdsws/php:3.0"
nodejs = "sdsws/node:2.1"
python2 =  "sdsws/python2:1.1"
python3 = "sdsws/python3:1.1"
golang = "sdsws/golang:1.1"
ruby = "sdsws/ruby:1.0"
mysql = "mysql:5.7"
mongodb = "mongo:4.2.1"
postgresql ="postgres:12.2-alpine"
```

You can replace the above default images and plug in your own docker images but make sure that each image has a **blocking CMD call** at the end of its corresponding dockerfile such as **CMD tail -f /dev/null**

For reference, you can check out the [dockerfiles](https://github.com/sdslabs/gasper-dockerfiles) for the default images used by Gasper
