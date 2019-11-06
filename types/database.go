package types

// DatabaseConfig is the configuration required for creating a database
type DatabaseConfig struct {
	Name     string `json:"name" valid:"required~Field 'name' is required but was not provided,alphanum~Field 'name' should only have alphanumeric characters,lowercase~Field 'name' should have only lowercase characters"`
	Password string `json:"password" valid:"required~Field 'password' is required but was not provided"`
}

// GetName returns the database's name
func (db *DatabaseConfig) GetName() string {
	return db.Name
}

// GetPassword returns the database's password
func (db *DatabaseConfig) GetPassword() string {
	return db.Password
}
