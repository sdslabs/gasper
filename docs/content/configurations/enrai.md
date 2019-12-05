# Enrai Configuration

Enrai service deals with reverse-proxying HTTP, HTTPS, HTTP/2, Websocket and gRPC requests to the desired application's IPv4 address and port based on the hostname

The following section deals with the configuration of Enrai

```toml
###########################
#   Enrai Configuration   #
###########################

[services.enrai]
# Time Interval (in seconds) in which `Enrai` updates its
# `Reverse-Proxy Record Storage` by polling the central registry-server.
record_update_interval = 30
deploy = false  # Deploy Enrai?
port = 80
```

!!!warning
    **Enrai** usually runs on port 80, hence the Gasper binary must be executed with **root** privileges in Linux

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

!!!warning
    **Enrai with SSL** usually runs on port 443, hence the Gasper binary must be executed with **root** privileges in Linux
