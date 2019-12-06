# Hikari Configuration

Hikari deals with creating and managing DNS records of all deployed applications

All DNS records point to the IPv4 addresses of Enrai ⚡ instances which in turn reverse-proxies the request to the desired application's IPv4 address and port

The following section deals with the configuration of Hikari

```toml
[services.hikari]
# Time Interval (in seconds) in which `Hikari` updates its
# `DNS Record Storage` by polling the central registry-server.
record_update_interval = 15
deploy = false  # Deploy Hikari?
port = 53
```

!!!warning
    You can reduce the value of **record_update_interval** parameter in the above configuration if you need changes in your ecosystem to propagate faster but this will in turn increase the load on the Redis central registry server so *choose wisely*

!!!warning
    **Hikari** usually runs on port 53, hence the Gasper binary must be executed with **root** privileges in Linux systems
