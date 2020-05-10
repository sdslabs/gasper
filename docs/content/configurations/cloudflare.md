# Cloudflare Configuration

Gasper has in-built support for using [cloudflare](https://www.cloudflare.com/) services

If Gasper's cloudflare plugin is enabled then whenever an application is created, its corresponding DNS entry will be automatically created in cloudflare

The DNS entry created will be according to the [domain](/configurations/global/#domain) parameter in the configuration file

???example
    If the domain parameter's value is `sdslabs.co` and you have created an application named **foo**, then an entry will be created in cloudflare (if plugin enabled) with the domain name `foo.app.sdslabs.co`

!!!warning
    The domain name set in the [domain](/configurations/global/#domain) parameter should be managed by cloudflare in order for this plugin to work

The following section deals with configurations related to Cloudflare

```toml
################################
#   CloudFlare Configuration   #
################################

[cloudflare]
# API Token used for creating/updating Cloudflare's DNS records.
# This token must have the scopes ZONE:ZONE:EDIT and ZONE:DNS:EDIT.
api_token = "" 
plugin = false  # Use Cloudflare Plugin?
public_ip = ""  # IPv4 address for Cloudflare's DNS records to point to.
```

You can generate a *Cloudflare API Token* from [here](https://dash.cloudflare.com/profile/api-tokens) and fill that value in the **api_token** field in the above configuration

!!!warning
    The generated token must have the permissions **ZONE:ZONE:EDIT** and **ZONE:DNS:EDIT** in order for this plugin to work

The **public_ip** field in the above configuration should hold the public IPv4 address of an **GenProxy ⚡** instance or a **load balancer** pointing to multiple **GenProxy ⚡** instances

!!!warning
    If you wish to use the Cloudflare plugin in your cloud ecosystem then make sure that the above configuration is **same** across all **nodes** where **AppMaker 💧** is deployed
