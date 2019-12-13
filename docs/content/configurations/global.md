# Global Configuration

This section of Gasper's configuration file deals with global settings

Here is the entire section of the configuration file dealing with global settings and we will go through each of them one by one

```toml
############################
#   Global Configuration   #
############################

# Run Gasper in Debug mode ?
# Set this value to `false` in Production.
debug = true

# Root domain for all deployed applications and databases.
domain = "sdslabs.co"

# Secret Key used for internal communication in the Gasper ecosystem.
secret = "YOUR_SECRET_KEY"

# Root of the deployed application in the docker container's filesystem.
project_root = "/gasper"

# Name of the file used for building and running applications.
# This file is application specific and must be present in an application's git repository's root.
# The contents of the file must be linux shell commands separated by newlines.
rc_file = "Gasperfile.txt"

# Run Gasper in Offline mode.
# For Development purposes only.
offline_mode = false

# DNS nameservers used by docker containers created by Gasper.
dns_servers = [
    "8.8.8.8",
    "8.8.4.4",
]
```

## Debug Mode

```toml
# Run Gasper in Debug mode ?
# Set this value to `false` in Production.
debug = true
```

The variable **debug** determines whether to run Gasper in debug mode or not 
In debug mode, internal server error messages are returned to the end user as JSON responses

!!!tip
    Set **debug** to `false` in Production 

## Domain

```toml
# Root domain for all deployed applications and databases.
domain = "sdslabs.co"
```

This section determines the root domain of all deployed applications and databases

The corresponding DNS entries for applications and databases will be automatically created by **Hikari 💡**

!!! example "DNS entry example for an application"
    If you create an application named **foo** then a DNS entry of `foo.app.sdslabs.co` will be created (based on the above root domain setting) pointing to the IPv4 address of an **Enrai ⚡** instance which in turn will reverse-proxy the request to the application's IPv4 address and port

!!! example "DNS entry example for a database"
    If you create a database named **bar** then a DNS entry of `bar.db.sdslabs.co` will be created (based on the above root domain setting) pointing to the IPv4 address of the node where the database's server is deployed

## Secret Key

```toml
# Secret Key used for internal communication in the Gasper ecosystem.
secret = "YOUR_SECRET_KEY"
```

Secret Key is used to encrypt the internal communications between **Kaze** 🌪 , **Mizu** 💧 and **Kaen** 🔥

!!!tip
    We recommend setting a strong **secret key** for securing your cloud ecosystem

!!!warning
    Make sure that the **secret key** is the same across all Gasper instances in your entire cloud ecosystem

## Project Root

```toml
# Root of the deployed application in the docker container's filesystem.
project_root = "/gasper"
```

All applications deployed by Gasper run within docker containers

**project_root** variable defines the root directory in the docker container's filesystem inside which the application's directory will be placed

## Run Commands File

```toml
# Name of the file used for building and running applications.
# This file is application specific and must be present in an application's git repository's root.
# The contents of the file must be linux shell commands separated by newlines.
rc_file = "Gasperfile.txt"
```

The Run Commands File or **rc_file** is the name of the file containing linux shell commands for building and running an application

This file must be present in an application's git repository's root directory

!!!info
    A user can deploy an application by either supplying the `build and run commands` in the request payload or by using this Run Commands File

???example "Usage"
    If the above **rc_file** parameter changes from `Gasperfile.txt` to `Alphadose`, then Gasper will look for a file named `Alphadose` in the application's git repository root during its deployment and will execute all commands present inside it

???example "Sample Run Commands File"
    For a [sample nodejs application](https://github.com/sdslabs/node), here is the corresponding run commands file [https://github.com/sdslabs/node/blob/master/Gasperfile.txt](https://github.com/sdslabs/node/blob/master/Gasperfile.txt)


## Offline Mode

```toml
# Run Gasper in Offline mode.
# For Development purposes only.
offline_mode = false
```

Gasper requires network connectivity for booting but with this parameter Gasper can run without it

Used for development purposes when the developer doesn't have an internet connectivity

!!!danger
    This functionality should be used strictly for development purposes

## DNS Nameservers

```toml
# DNS nameservers used by docker containers created by Gasper.
dns_servers = [
    "8.8.8.8",
    "8.8.4.4",
]
```

This field defines the DNS Nameservers that would be used inside all deployed application's docker containers for domain name resolution

Change it according to your network infrastructure if required

!!!info
    By default Google's nameservers are used
