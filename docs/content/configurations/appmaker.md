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
```

!!!warning
    The node where **AppMaker** is to be deployed should have **Docker** installed and running
