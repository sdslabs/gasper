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
	mysqlRootPassword = configs.ServiceConfig.Kaen.MySQL.Env["MYSQL_ROOT_PASSWORD"].(string)
)

// CreateMysqlDB creates a database in the Mysql instance with the given database name, user and password
func CreateMysqlDB(db types.Database) error {
	port := configs.ServiceConfig.Kaen.MySQL.ContainerPort

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%d)", port)
	connection := fmt.Sprintf("%s:%s@%s/", mysqlRootUser, mysqlRootPassword, agentAddress)

	conn, err := sql.Open(mysqlDriver, connection)

	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}
	defer conn.Close()

	_, err = conn.Exec("CREATE DATABASE " + db.GetName())
	if err != nil {
		return fmt.Errorf("Error while creating the database : Database Already Exists")
	}

	query := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", db.GetUser(), mysqlHost, db.GetPassword())
	_, err = conn.Exec(query)
	if err != nil {
		errs := refreshDBUser(db, conn)
		if errs != nil {
			return fmt.Errorf("Error while creating the database : %s", err)
		}
	}

	query = fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%s'", db.GetName(), db.GetUser(), mysqlHost)
	_, err = conn.Exec(query)
	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}

	_, err = conn.Exec("FLUSH PRIVILEGES")
	if err != nil {
		return fmt.Errorf("Error while flushing user priviliges : %s", err)
	}

	return nil
}

// DeleteMysqlDB deletes the database given by the database name and username
func DeleteMysqlDB(databaseName string) error {
	username := databaseName
	port := configs.ServiceConfig.Kaen.MySQL.ContainerPort

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%d)", port)
	connection := fmt.Sprintf("%s:%s@%s/", mysqlRootUser, mysqlRootPassword, agentAddress)

	conn, err := sql.Open(mysqlDriver, connection)
	if err != nil {
		return fmt.Errorf("Error while connecting to database : %s", err)
	}
	defer conn.Close()

	_, err = conn.Exec("DROP DATABASE IF EXISTS " + databaseName)
	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}

	_, err = conn.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'@'%s'", username, mysqlHost))
	if err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err)
	}
	return nil
}

func refreshDBUser(db types.Database, conn *sql.DB) error {
	_, errf := conn.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'@'%s'", db.GetUser(), mysqlHost))
	if errf != nil {
		return fmt.Errorf("Error while deleting the user : %s", errf)
	}

	query := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", db.GetUser(), mysqlHost, db.GetPassword())
	_, errc := conn.Exec(query)
	if errc != nil {
		return fmt.Errorf("Error while creating the database : %s", errc)
	}

	return nil
}
