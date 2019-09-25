package mongo

import "go.mongodb.org/mongo-driver/mongo"

const (
	// DBInstance is db instance type name in the instances collection
	DBInstance = "db"
	// Mysql is db instance type name for mysql database in the instances collection
	Mysql = "mysql"

	// MongoDB is db instance type name in the instances collection
	MongoDB = "mongodb"

	// AppInstance is app instance type name in the instances collection
	AppInstance = "app"

	// InstanceCollection is the collection for all the instances
	InstanceCollection = "instances"

	// UserCollection is the collection for all users
	UserCollection = "users"
)

// ErrNoDocuments is the error when no matching documents are found
// for an update operation
var ErrNoDocuments = mongo.ErrNoDocuments
