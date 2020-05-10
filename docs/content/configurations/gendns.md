# GenDNS Configuration

GenDNS deals with creating and managing DNS records of all deployed applications and databases

All application DNS records point to the IPv4 addresses of GenProxy ⚡ instances which in turn reverse-proxies the request to the desired application's IPv4 address and port

All database DNS records point to the IPv4 address of the node where the database's server is deployed

!!!info
    **GenDNS 💡** automatically creates a DNS entry for **Master 🌪** (if deployed) pointing to an **GenProxy ⚡** instance which will be further load-balanced among all available **Master 🌪** instances

    The created DNS entry will be based on the [domain](/configurations/global/#domain) parameter

    !!!example
        If the [domain](/configurations/global/#domain) parameter is set to `sdslabs.co` then the corresponding DNS entry `master.sdslabs.co` will be created by **GenDNS 💡**

The following section deals with the configuration of GenDNS

```toml
############################
#   GenDNS Configuration   #
############################

[services.gendns]
# Time Interval (in seconds) in which `GenDNS` updates its
# `DNS Record Storage` by polling the central registry-server.
record_update_interval = 15
deploy = false  # Deploy GenDNS?
port = 53
```

!!!tip
    You can reduce the value of **record_update_interval** parameter in the above configuration if you need changes in your ecosystem to propagate faster but this will in turn increase the load on the Redis central registry server so *choose wisely*

!!!warning
    **GenDNS** usually runs on port 53, hence the Gasper binary must be executed with **root** privileges in Linux systems
