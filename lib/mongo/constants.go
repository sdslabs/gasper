package mongo

import (
	"github.com/sdslabs/gasper/types"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// projectDatabase is the name of the database used for storing all of gasper's information
	projectDatabase = "gasper"

	// DBInstance is db instance type name in the instances collection
	DBInstance = "database"

	// Mysql is db instance type name for mysql database in the instances collection
	Mysql = types.MySQL

	// MongoDB is db instance type name in the instances collection
	MongoDB = types.MongoDB

	// AppInstance is app instance type name in the instances collection
	AppInstance = "application"

	// InstanceCollection is the collection for all the instances
	InstanceCollection = "instances"

	// UserCollection is the collection for all users
	UserCollection = "users"

	// MetricsCollection is the collection to hold the metrics of the instances
	MetricsCollection = "metrics"

	// NameKey is the key holding the name of an instance
	NameKey = "name"

	// OwnerKey is the key holding the owner of an instance
	OwnerKey = "owner"

	// InstanceTypeKey is the key holding the instance type of an instance
	InstanceTypeKey = "instance_type"

	// LanguageKey is the key holding the language of an instance
	LanguageKey = "language"

	// HostIPKey is the key holding the host IP address of an instance
	HostIPKey = "host_ip"

	// ContainerPortKey is the key holding the port of the container in which an application is deployed
	ContainerPortKey = "container_port"

	// PortKey is the key holding the port of the container in which a database server is deployed
	PortKey = "port"

	// EmailKey is the key holding the email of a user
	EmailKey = "email"

	// UsernameKey is the key holding the username of a user
	UsernameKey = "username"

	// PasswordKey is the key holding the password of a user/instance
	PasswordKey = "password"

	// AdminKey is the key denoting whether a user has superuser privileges or not
	AdminKey = "admin"

	// TimestampKey is the key holding the timestamp of when a metrics collection was inserted
	TimestampKey = "timestamp"
)

// ErrNoDocuments is the error when no matching documents are found
// for an update operation
var ErrNoDocuments = mongo.ErrNoDocuments
