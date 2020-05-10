# Creating a MySQL Database

This example shows how to deploy a [MySQL](https://www.mysql.com/) database via Gasper

!!!warning "Prerequisites"
    * You have [Master](/configurations/master/) and [DbMaker](/configurations/dbmaker/) up and running
    * You have [DbMaker MySQL Plugin](/configurations/dbmaker/#mysql-configuration) enabled
    * You have already [logged in](/examples/login/) and obtained a JSON Web Token

```bash
$ curl -X POST \
  http://localhost:3000/dbs/mysql \
  -H 'Authorization: Bearer {{token}}' \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "alphamysql",
	"password": "alphamysql"
}'

{
    "name": "alphamysql",
    "password": "alphamysql",
    "user": "alphamysql",
    "instance_type": "database",
    "language": "mysql",
    "host_ip": "10.43.3.24",
    "port": 33061,
    "owner": "anish.mukherjee1996@gmail.com",
    "success": true
}
```
