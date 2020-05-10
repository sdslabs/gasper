# JWT Configuration

This is used to configure the timeout and refresh time for the authentication token

Both values `timeout` and `max_refresh` are defined in the **jwt** section of the configuration file

`timeout` is the interval after which the client needs to request for a new token which can be done by either logging in again or obtaining the token through refresh route (`GET /auth/refresh`)

`max_refresh` is the time interval after which user is logged out and the token can only be obtained by logging in again

```toml
#########################
#   JWT Configuration   #
#########################

# Configuration for the JSON Web Token (JWT) authentication mechanism.
[jwt]

# timeout refers to the duration in which the JWT is valid.
# max_refresh refers to the duration in which the JWT can be refreshed after its expiry.

# Both timeout and max_refresh are in seconds
# Total refresh time = max_refresh + timeout
timeout = 3600 # 1 hour
max_refresh = 2419200 # 28 days
```

!!!info
    The above section only needs to be configured for the nodes where **Master 🌪** is to be deployed
