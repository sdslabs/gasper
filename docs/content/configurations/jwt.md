# JWT Configuration

This is used to configure the timeout and refresh time for the authentication token.

Both values `timeout` and `max_refresh` are defined in the **jwt** section of the configuration file.

`timeout` is the interval after which the client needs to request for a new token. The new token can be requested by either logging in again or obtaining the token through refresh route (`GET /auth/refresh`).

`max_refresh` is the time interval after which user is logged out and the token can only be obtained by logging in again.

```toml
#########################
#   JWT Configuration   #
#########################

# Both timeout and max_refresh in seconds
# Total refresh time = max_refresh + timeout
# Since, max_refresh >> timeout
# total refresh ~ max_refresh
[jwt]
timeout = 3600 # 1 hour
max_refresh = 2419200 # 28 days
```
