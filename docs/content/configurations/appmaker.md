# AppMaker Configuration

AppMaker deals with creating and managing applications and their life-cycles

The following section deals with the configuration of AppMaker

```toml
##############################
#   AppMaker Configuration   #
##############################

[services.appmaker]
deploy = true   # Deploy AppMaker?
port = 4000

# Time Interval (in seconds) in which metrics of all application containers
# running in the current node are collected and stored in the central mongoDB database
metrics_interval = 600

# Time Interval (in seconds) in which health is checked of all application containers and if unhealthy, they are restarted
health_interval = 300

# Hard Limits the total number of app instances that can be deployed by an user
# Set app_limit = -1 if no hard limit is to be imposed
app_limit = 10

# Specifies the maximum CPU allocation for a container created by a non admin user
max_container_cpu = 0.25

# Specifies the maximum memory allocation for a container created by a non admin user
max_container_memory = 0.5
```

!!!warning
    The node where **AppMaker** is to be deployed should have **Docker** installed and running
