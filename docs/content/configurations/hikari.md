# Hikari Configuration

Hikari deals with creating and managing DNS records of all deployed applications and databases

All application DNS records point to the IPv4 addresses of Enrai ⚡ instances which in turn reverse-proxies the request to the desired application's IPv4 address and port

All database DNS records point to the IPv4 address of the node where the database's server is deployed

!!!info
    **Hikari 💡** automatically creates a DNS entry for **Kaze 🌪** (if deployed) pointing to an **Enrai ⚡** instance which will be further load-balanced among all available **Kaze 🌪** instances

The following section deals with the configuration of Hikari

```toml
############################
#   Hikari Configuration   #
############################

[services.hikari]
# Time Interval (in seconds) in which `Hikari` updates its
# `DNS Record Storage` by polling the central registry-server.
record_update_interval = 15
deploy = false  # Deploy Hikari?
port = 53
```

!!!tip
    You can reduce the value of **record_update_interval** parameter in the above configuration if you need changes in your ecosystem to propagate faster but this will in turn increase the load on the Redis central registry server so *choose wisely*

!!!warning
    **Hikari** usually runs on port 53, hence the Gasper binary must be executed with **root** privileges in Linux systems
