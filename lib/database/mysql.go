package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/sdslabs/gasper/configs"
)

var (
	mysqlDriver       = "mysql"
	mysqlHost         = `%`
	mysqlRootUser     = "root"
	mysqlRootPassword = configs.ServiceConfig.Mysql.Env["MYSQL_ROOT_PASSWORD"].(string)
)

type mysqlAgentServer struct{}

// CreateMysqlDB creates a database in the Mysql instance with the given database name, user and password
func CreateMysqlDB(database, username, password string) error {
	port := configs.ServiceConfig.Mysql.ContainerPort

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%d)", port)
	connection := fmt.Sprintf("%s:%s@%s/", mysqlRootUser, mysqlRootPassword, agentAddress)

	db, err := sql.Open(mysqlDriver, connection)

	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + database)
	if err != nil {
		return fmt.Errorf("Error while creating the database : Database Already Exists")
	}

	query := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", username, mysqlHost, password)
	_, err = db.Exec(query)
	if err != nil {
		errs := refreshDBUser(database, username, password, db)
		if errs != nil {
			return fmt.Errorf("Error while creating the database : %s", err)
		}
	}

	query = fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%s'", database, username, mysqlHost)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}

	_, err = db.Exec("FLUSH PRIVILEGES")
	if err != nil {
		return fmt.Errorf("Error while flushing user priviliges : %s", err)
	}

	return nil
}

// DeleteMysqlDB deletes the database given by the database name and username
func DeleteMysqlDB(database string) error {
	username := database
	port := configs.ServiceConfig.Mysql.ContainerPort

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%d)", port)
	connection := fmt.Sprintf("%s:%s@%s/", mysqlRootUser, mysqlRootPassword, agentAddress)

	db, err := sql.Open(mysqlDriver, connection)
	if err != nil {
		return fmt.Errorf("Error while connecting to database : %s", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS " + database)
	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}

	_, err = db.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'@'%s'", username, mysqlHost))
	if err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err)
	}
	return nil
}

func refreshDBUser(database, username, password string, db *sql.DB) error {
	_, errf := db.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'@'%s'", username, mysqlHost))
	if errf != nil {
		return fmt.Errorf("Error while deleting the user : %s", errf)
	}

	query := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", username, mysqlHost, password)
	_, errc := db.Exec(query)
	if errc != nil {
		return fmt.Errorf("Error while creating the database : %s", errc)
	}

	return nil
}
