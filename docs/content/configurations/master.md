# Master Configuration

Master is the master of the entire Gasper ecosystem which performs the following tasks

* Equal distribution of applications and databases among worker nodes
* User Authentication based on JWT (JSON Web Token)
* User API for performing operations on any application/database in any node (Identity Access Management is handled with JWT)
* Admin API for fetching and managing information of all nodes, applications, databases and users
* Removal of inactive nodes from the cloud ecosystem
* Re-scheduling of applications in case of node failure

Master API docs are available [here](/api)

The following section deals with the configuration of Master

```toml
############################
#   Master Configuration   #
############################

[services.master]
# Time Interval (in seconds) in which `Master` sends health-check probes
# to all worker nodes and removes inactive nodes from the central registry-server.
cleanup_interval = 600
deploy = true   # Deploy Master?
port = 3000

# Configuration for the MongoDB service container required by all deployed services.
[services.master.mongodb]
plugin = true  # Deploy MongoDB server and let `Master` manage it?
container_port = 27019  # Port on which the MongoDB server container will run

# Environment variables for MongoDB docker container.
[services.master.mongodb.env]
MONGO_INITDB_ROOT_USERNAME = "alphadose"   # Root user of MongoDB server inside the container
MONGO_INITDB_ROOT_PASSWORD = "alphadose"   # Root password of MongoDB server inside the container

# Configuration for the Redis service container required by all deployed services.
[services.master.redis]
plugin = true  # Deploy Redis server and let `Master` manage it?
container_port = 6380  # Port on which the Redis server container will run
password = "alphadose"
```

!!!tip
    You can reduce the value of **cleanup_interval** parameter in the above configuration if you need changes in your ecosystem to propagate faster but this will in turn increase the load on the Redis central registry server so *choose wisely*
