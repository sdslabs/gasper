# Creating a MongoDB Database

This example shows how to deploy a [MongoDB](https://www.mongodb.com/) database via Gasper

!!!warning "Prerequisites"
    * You have [Kaze](/configurations/kaze/) and [Kaen](/configurations/kaen/) up and running
    * You have [Kaen MongoDB Plugin](/configurations/kaen/#mongodb-configuration) enabled
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token

```bash
$ curl -X POST \
  http://localhost:3000/dbs/mongodb \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "alphamongo",
	"password": "alphamongo"
}'

{
    "name": "alphamongo",
    "password": "alphamongo",
    "user": "alphamongo",
    "instance_type": "database",
    "language": "mongodb",
    "host_ip": "10.43.3.24",
    "port": 27018,
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```
