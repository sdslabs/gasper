# MongoDB Configuration

Gasper uses [MongoDB](https://www.mongodb.com/) for storing data pertaining to users, applications and databases

The following section in the configuration file deals with MongoDB

```toml
#############################
#   MongoDB Configuration   #
#############################

[mongo]
# For databases with authentication
# use the following URL format `mongodb://username:password@host:port`.
url = "mongodb://alphadose:alphadose@localhost:27019/?authSource=admin"
```

!!!warning
    There should be only a single instance of MongoDB running in your entire cloud ecosystem and all instances of Gasper should connect only to that single MongoDB instance i.e the above configuration must be **same** across all Gasper instances in all nodes
