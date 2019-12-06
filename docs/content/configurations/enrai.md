# Enrai Configuration

Enrai service deals with reverse-proxying HTTP, HTTPS, HTTP/2, Websocket and gRPC requests to the desired application's IPv4 address and port based on the hostname

!!!info
    **Enrai ⚡** automatically creates a reverse-proxy entry for **Kaze 🌪** (if deployed) pointing to its IPv4 address and port

## Default
The following section deals with the configuration of Enrai

```toml
###########################
#   Enrai Configuration   #
###########################

[services.enrai]
# Time Interval (in seconds) in which `Enrai` updates its
# `Reverse-Proxy Record Storage` by polling the central registry-server.
record_update_interval = 15
deploy = false  # Deploy Enrai?
port = 80
```

!!!tip
    You can reduce the value of **record_update_interval** parameter in the above configuration if you need changes in your ecosystem to propagate faster but this will in turn increase the load on the Redis central registry server so *choose wisely*

!!!warning
    **Enrai** usually runs on port 80, hence the Gasper binary must be executed with **root** privileges in Linux systems

## Enrai with SSL

The following section deals with configuring Enrai with SSL support for HTTPS

```toml
# Configuration for using SSL with `Enrai`.
[services.enrai.ssl]
plugin = false  # Use SSL with Enrai?
port = 443
certificate = "/home/user/fullchain.pem"  # Certificate Location
private_key = "/home/user/privkey.pem"  # Private Key Location
```

The **certificate** and **private key** in the above configuration should be configured for all sub-domains based on the [domain parameter](/configurations/global/#domain) in the configuration file

!!!example "Configuration Example"
    If the **domain** parameter is **sdslabs.co** then the certificate and private key should be configured for the following subdomains `*.sdslabs.co` and `*.*.sdslabs.co`

!!!warning
    **Enrai with SSL** usually runs on port 443, hence the Gasper binary must be executed with **root** privileges in Linux systems
