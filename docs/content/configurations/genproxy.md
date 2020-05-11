# GenProxy Configuration

GenProxy service deals with reverse-proxying HTTP, HTTPS, HTTP/2, Websocket and gRPC requests to the desired application's IPv4 address and port based on the hostname

!!!info
    **GenProxy ⚡** automatically creates a reverse-proxy entry for **Master 🌪** (if deployed) pointing to its IPv4 address and port

## Default
The following section deals with the configuration of GenProxy

```toml
##############################
#   GenProxy Configuration   #
##############################

[services.genproxy]
# Time Interval (in seconds) in which `GenProxy` updates its
# `Reverse-Proxy Record Storage` by polling the central registry-server.
record_update_interval = 15
deploy = false  # Deploy GenProxy?
port = 80
```

!!!tip
    You can reduce the value of **record_update_interval** parameter in the above configuration if you need changes in your ecosystem to propagate faster but this will in turn increase the load on the Redis central registry server so *choose wisely*

!!!warning
    **GenProxy** usually runs on port 80, hence the Gasper binary must be executed with **root** privileges in Linux systems

## GenProxy with SSL

The following section deals with configuring GenProxy with SSL support for HTTPS

```toml
# Configuration for using SSL with `GenProxy`.
[services.genproxy.ssl]
plugin = false  # Use SSL with GenProxy?
port = 443
certificate = "/home/user/fullchain.pem"  # Certificate Location
private_key = "/home/user/privkey.pem"  # Private Key Location
```

The **certificate** and **private key** in the above configuration should be configured for all sub-domains based on the [domain parameter](/configurations/global/#domain) in the configuration file

!!!example "Configuration Example"
    If the [domain](/configurations/global/#domain) parameter is `sdslabs.co` then the certificate and private key should be configured for the following subdomains `*.sdslabs.co` and `*.*.sdslabs.co`

!!!warning
    **GenProxy with SSL** usually runs on port 443, hence the Gasper binary must be executed with **root** privileges in Linux systems
