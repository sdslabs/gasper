package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sdslabs/SWS/lib/utils"
)

var dbHost = `%`
var dbUser = "root"

type mysqlAgentServer struct{}

// CreateDB creates a database in the Mysql instance with the given database name, user and password
func CreateDB(database, username, password string) error {
	port := utils.ServiceConfig["mysql"].(map[string]interface{})["container_port"].(string)

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)
	connection := fmt.Sprintf("%s@%s/", dbUser, agentAddress)

	db, err := sql.Open("mysql", connection)

	if err != nil {
		return fmt.Errorf("Error while creating the database : %s", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + database)
	if err != nil {
		errs := sanitaryActions(database, username, password, db, 1)
		if errs != nil {
			return fmt.Errorf("Error while creating the database : %s", err)
		}
	}

	query := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", username, dbHost, password)
	_, err = db.Exec(query)
	if err != nil {
		errs := sanitaryActions(database, username, password, db, 2)
		if errs != nil {
			return fmt.Errorf("Error while creating the database : %s", err)
		}
	}

	query = fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%s'", database, username, dbHost)
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

// DeleteDB deletes the database given by the database name and username
func DeleteDB(database, username string) error {
	port := utils.ServiceConfig["mysql"].(map[string]interface{})["container_port"].(string)

	agentAddress := fmt.Sprintf("tcp(127.0.0.1:%s)", port)
	connection := fmt.Sprintf("%s@%s/", dbUser, agentAddress)

	db, err := sql.Open("mysql", connection)
	if err != nil {
		return fmt.Errorf("Error while connecting to database : %s", err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE " + database)
	if err != nil {
		return fmt.Errorf("Error while deleting the database : %s", err)
	}

	_, err = db.Exec(fmt.Sprintf("DROP USER '%s'@'%s'", username, dbHost))
	if err != nil {
		return fmt.Errorf("Error while deleting the user : %s", err)
	}
	return nil
}

func sanitaryActions(database, username, password string, db *sql.DB, stage int) error {
	switch stage {
	case 1:
		{
			_, errf := db.Exec("DROP DATABASE " + database)
			if errf != nil {
				return fmt.Errorf("Error while deleting the database : %s", errf)
			}

			_, errc := db.Exec("CREATE DATABASE " + database)
			if errc != nil {
				return fmt.Errorf("Error while creating the database : %s", errc)
			}
		}

	case 2:
		{
			_, errf := db.Exec(fmt.Sprintf("DROP USER '%s'@'%s'", username, dbHost))
			if errf != nil {
				return fmt.Errorf("Error while deleting the user : %s", errf)
			}

			query := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", username, dbHost, password)
			_, errc := db.Exec(query)
			if errc != nil {
				return fmt.Errorf("Error while creating the database : %s", errc)
			}
		}
	}

	return nil
}
