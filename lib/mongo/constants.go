package mongo

import (
	"github.com/sdslabs/gasper/types"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
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

	// InstanceTypeKey is the key holding the instance type of an instance
	InstanceTypeKey = "instance_type"

	// HostIPKey is the key holding the host IP address of an instance
	HostIPKey = "host_ip"

	// ContainerPortKey is the key holding the port of the container in which an application is deployed
	ContainerPortKey = "container_port"
)

// ErrNoDocuments is the error when no matching documents are found
// for an update operation
var ErrNoDocuments = mongo.ErrNoDocuments
