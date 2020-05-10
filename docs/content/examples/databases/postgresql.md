# Creating a PostgreSQL Database

This example shows how to deploy a [PostgreSQL](https://www.postgresql.org/) database via Gasper

!!!warning "Prerequisites"
    * You have [Master](/configurations/master/) and [DbMaker](/configurations/dbmaker/) up and running
    * You have [DbMaker PostgreSQL Plugin](/configurations/dbmaker/#postgresql-configuration) enabled
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token

```bash
$ curl -X POST \
  http://localhost:3000/dbs/postgresql \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "alphapostgresql",
	"password": "alphapostgresql"
}'

{
    "name": "alphapostgresql",
    "password": "alphapostgresql",
    "user": "alphapostgresql",
    "instance_type": "database",
    "language": "postgresql",
    "db_url": "alphapostgresql.db.sdslabs.co",
    "host_ip": "192.168.225.90",
    "port": 29121,
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```
