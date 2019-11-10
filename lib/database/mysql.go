package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/types"
)

var (
	mysqlDriver       = "mysql"
	mysqlHost         = `%`
	mysqlRootUser     = "root"
	mysqlPort         = configs.ServiceConfig.Kaen.MySQL.ContainerPort
	mysqlRootPassword = configs.ServiceConfig.Kaen.MySQL.Env["MYSQL_ROOT_PASSWORD"].(string)
)

// CreateMysqlDB creates a database in the Mysql instance with the given database name, user and password
func CreateMysqlDB(db types.Database) error {
	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%d)", mysqlPort)
	connection := fmt.Sprintf("%s:%s@%s/", mysqlRootUser, mysqlRootPassword, agentAddress)

	conn, err := sql.Open(mysqlDriver, connection)
	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}
	defer conn.Close()

	if _, err = conn.Exec("CREATE DATABASE " + db.GetName()); err != nil {
		return fmt.Errorf("Error while creating the database : Database Already Exists")
	}

	query := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", db.GetUser(), mysqlHost, db.GetPassword())
	if _, err = conn.Exec(query); err != nil {
		if err = refreshDBUser(db, conn); err != nil {
			return fmt.Errorf("Error while creating the database : %s", err)
		}
	}

	query = fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%s'", db.GetName(), db.GetUser(), mysqlHost)
	if _, err = conn.Exec(query); err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}

	if _, err = conn.Exec("FLUSH PRIVILEGES"); err != nil {
		return fmt.Errorf("Error while flushing user priviliges : %s", err)
	}

	return nil
}

// DeleteMysqlDB deletes the database given by the database name and username
func DeleteMysqlDB(databaseName string) error {
	username := databaseName

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%d)", mysqlPort)
	connection := fmt.Sprintf("%s:%s@%s/", mysqlRootUser, mysqlRootPassword, agentAddress)

	conn, err := sql.Open(mysqlDriver, connection)
	if err != nil {
		return fmt.Errorf("Error while connecting to database : %s", err)
	}
	defer conn.Close()

	if _, err = conn.Exec("DROP DATABASE IF EXISTS " + databaseName); err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}

	if _, err = conn.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'@'%s'", username, mysqlHost)); err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err)
	}
	return nil
}

func refreshDBUser(db types.Database, conn *sql.DB) error {
	_, err := conn.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'@'%s'", db.GetUser(), mysqlHost))
	if err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err)
	}

	query := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", db.GetUser(), mysqlHost, db.GetPassword())
	if _, err = conn.Exec(query); err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}
	return nil
}
