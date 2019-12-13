# Redis Configuration

Gasper uses [Redis](https://redis.io/) as a central registry server

It is used for storing the addresses of all active components, nodes, applications and databases

```toml
###########################
#   Redis Configuration   #
###########################

# Acts as a central-registry for the Gasper ecosystem
[redis]
host = "localhost"
port = 6379
password = ""
db = 0
```

!!!warning
    There should be only a single instance of Redis running in your entire cloud ecosystem and all instances of Gasper should connect only to that single Redis instance i.e the above configuration must be **same** across all Gasper instances in all nodes
