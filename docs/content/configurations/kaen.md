# Kaen Configuration

Kaen deals with creating and managing databases and their life-cycles

The following section deals with the configuration of Kaen

```toml
##########################
#   Kaen Configuration   #
##########################

[services.kaen]
deploy = false   # Deploy Kaen?
port = 9000
```

!!!info
    The node where **Kaen** is to be deployed should have **Docker** installed and running

## MySQL Configuration

This section deals with the MySQL server configuration managed by Kaen

```toml
# Configuration for MySQL database server managed by `Kaen`
[services.kaen.mysql]
plugin = false  # Deploy MySQL server and let `Kaen` manage it?
container_port = 33061  # Port on which the MySQL server container will run

# Environment variables for MySQL docker container.
[services.kaen.mysql.env]
MYSQL_ROOT_PASSWORD = "YOUR_MYSQL_PASSWORD"  # Root password of MySQL server inside the container
```

!!!info
    The username of the deployed MySQL server will be **root** and the password will be the value of the variable **MYSQL_ROOT_PASSWORD**

## MongoDB Configuration

This section deals with the MongoDB server configuration managed by Kaen

```toml
# Configuration for MongoDB database server managed by `Kaen`
[services.kaen.mongodb]
plugin = false  # Deploy MongoDB server and let `Kaen` manage it?
container_port = 27018  # Port on which the MongoDB server container will run

# Environment variables for MongoDB docker container.
[services.kaen.mongodb.env]
MONGO_INITDB_ROOT_USERNAME = "YOUR_ROOT_NAME"   # Root user of MongoDB server inside the container
MONGO_INITDB_ROOT_PASSWORD = "YOUR_ROOT_PASSWORD"   # Root password of MongoDB server inside the container
```

!!!info
    The username of the deployed MongoDB server will be the value of the variable **MONGO_INITDB_ROOT_USERNAME** and the password will be the value of the variable **MONGO_INITDB_ROOT_PASSWORD**
