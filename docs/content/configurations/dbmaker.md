# DbMaker Configuration

DbMaker deals with creating and managing databases and their life-cycles

The following section deals with the configuration of DbMaker

```toml
##########################
#   DbMaker Configuration   #
##########################

[services.dbmaker]
deploy = false   # Deploy DbMaker?
port = 9000
```

!!!warning
    The node where **DbMaker** is to be deployed should have **Docker** installed and running

## MySQL Configuration

This section deals with the MySQL server configuration managed by DbMaker

```toml
# Configuration for MySQL database server managed by `DbMaker`
[services.dbmaker.mysql]
plugin = false  # Deploy MySQL server and let `DbMaker` manage it?
container_port = 33061  # Port on which the MySQL server container will run

# Environment variables for MySQL docker container.
[services.dbmaker.mysql.env]
MYSQL_ROOT_PASSWORD = "YOUR_MYSQL_PASSWORD"  # Root password of MySQL server inside the container
```

!!!info
    The username of the deployed MySQL server will be **root** and the password will be the value of the variable **MYSQL_ROOT_PASSWORD**

## MongoDB Configuration

This section deals with the MongoDB server configuration managed by DbMaker

```toml
# Configuration for MongoDB database server managed by `DbMaker`
[services.dbmaker.mongodb]
plugin = false  # Deploy MongoDB server and let `DbMaker` manage it?
container_port = 27018  # Port on which the MongoDB server container will run

# Environment variables for MongoDB docker container.
[services.dbmaker.mongodb.env]
MONGO_INITDB_ROOT_USERNAME = "YOUR_ROOT_NAME"   # Root user of MongoDB server inside the container
MONGO_INITDB_ROOT_PASSWORD = "YOUR_ROOT_PASSWORD"   # Root password of MongoDB server inside the container
```

!!!info
    The username of the deployed MongoDB server will be the value of the variable **MONGO_INITDB_ROOT_USERNAME** and the password will be the value of the variable **MONGO_INITDB_ROOT_PASSWORD**

## PostgreSQL Configuration

This section deals with the PostgreSQL server configuration managed by DbMaker

```toml
# Configuration for PostgreSQL database server managed by `DbMaker`
[services.dbmaker.postgresql]
plugin = false  # Deploy PostgreSQL server and let `DbMaker` manage it?
container_port = 29121  # Port on which the PostgreSQL server container will run

# Environment variables for PostgreSQL docker container.
[services.dbmaker.postgresql.env]
POSTGRES_USER = "YOUR_ROOT_NAME"   # Root user of PostgreSQL server inside the container
POSTGRES_PASSWORD = "YOUR_ROOT_PASSWORD"   # Root password of PostgreSQL server inside the container
```

!!!info
    The username of the deployed PostgreSQL server will be the value of the variable **POSTGRES_USER** and the password will be the value of the variable **POSTGRES_PASSWORD**

## Redis Configuration

This section deals with the Redis server configuration managed by DbMaker

```toml
# Configuration for Redis database server managed by `DbMaker`
[services.dbmaker.redis]
plugin = false  # Deploy RedisDB server and let `DbMaker` manage it

```

!!!info
    * For Redis due to the lack of namespaces a new container is created per user unlike others where one database is created per user in a single container
    * The container name of the deployed Redis server will be the value of the variable **username** and the password will be the value of the variable **password** both of which are retrieved from the API request to the master service
