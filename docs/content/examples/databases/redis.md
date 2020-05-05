# Creating a Redis Database

This example shows how to deploy a [Redis](https://redis.io/) database via Gasper

!!!warning "Prerequisites"
    * You have [Kaze](/configurations/kaze/) and [Kaen](/configurations/kaen/) up and running
    * You have [Kaen Redis Plugin](/configurations/kaen/#redis-configuration) enabled
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token

```bash
$ curl -X POST \
  http://localhost:3000/dbs/redis \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "alpha",
	"password": "alpha"
}'

{
    "name": "alpha",
    "password": "alpha",
    "user": "alpha",
    "instance_type": "database",
    "language": "redis",
    "db_url": "alpha3.db.sdslabs.co",
    "host_ip": "192.168.43.46",
    "port": 45861,
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```
