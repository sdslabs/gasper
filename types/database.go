package types

// Database is the interface for creating a database
type Database interface {
	GetName() string
	GetPassword() string
	GetUser() string
}

// DatabaseConfig is the configuration required for creating a database
type DatabaseConfig struct {
	Name          string `json:"name" bson:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,lowercase~Field 'name' should have only lowercase characters"`
	Password      string `json:"password" bson:"password" valid:"required~Field 'password' is required but was not provided"`
	User          string `json:"user,omitempty" bson:"user,omitempty"`
	InstanceType  string `json:"instance_type,omitempty" bson:"instance_type,omitempty"`
	Language      string `json:"language,omitempty" bson:"language,omitempty"`
	HostIP        string `json:"host_ip,omitempty" bson:"host_ip,omitempty"`
	ContainerPort int    `json:"port,omitempty" bson:"port,omitempty"`
	Owner         string `json:"owner,omitempty" bson:"owner,omitempty"`
	Success       bool   `json:"success,omitempty" bson:"-"`
}

// GetName returns the database's name
func (db *DatabaseConfig) GetName() string {
	return db.Name
}

// GetPassword returns the database's password
func (db *DatabaseConfig) GetPassword() string {
	return db.Password
}

// SetUser sets the database's user in its context
func (db *DatabaseConfig) SetUser(user string) {
	db.User = user
}

// GetUser returns the database's user
func (db *DatabaseConfig) GetUser() string {
	if db.User == "" {
		return db.Name
	}
	return db.User
}

// SetInstanceType sets the database's type of instance in its context
func (db *DatabaseConfig) SetInstanceType(instanceType string) {
	db.InstanceType = instanceType
}

// SetLanguage sets the database's language in its context
func (db *DatabaseConfig) SetLanguage(language string) {
	db.Language = language
}

// SetHostIP sets the IP address of the host in which the database is deployed
// in its context
func (db *DatabaseConfig) SetHostIP(IP string) {
	db.HostIP = IP
}

// SetContainerPort sets the port in which the database server is running
// in the host system to the database's context
func (db *DatabaseConfig) SetContainerPort(port int) {
	db.ContainerPort = port
}

// SetOwner sets the owner of the database in its context
// The owner is referenced by his/her email ID
func (db *DatabaseConfig) SetOwner(owner string) {
	db.Owner = owner
}

// SetSuccess defines the success of creating the database
func (db *DatabaseConfig) SetSuccess(success bool) {
	db.Success = success
}
