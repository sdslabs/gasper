# Mizu Configuration

Mizu deals with creating and managing applications and their life-cycles

The following section deals with the configuration of Mizu

```toml
##########################
#   Mizu Configuration   #
##########################

[services.mizu]
deploy = true   # Deploy Mizu?
port = 4000
# Time Interval (in seconds) in which metrics of all application containers
# running in the current node are collected and stored in the central mongoDB database
metrics_interval = 600
```

!!!warning
    The node where **Mizu** is to be deployed should have **Docker** installed and running
