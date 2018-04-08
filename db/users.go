package db

import (
	"database/sql"
	"bytes"
	"fmt"
	"log"
)

// Creates the users table
func SetupUsersTable(db *sql.DB) {
	var createTableCommand bytes.Buffer
	createTableCommand.WriteString("CREATE TABLE IF NOT EXISTS ")
	createTableCommand.WriteString(userTableName)
	createTableCommand.WriteString(" (\n")
	createTableCommand.WriteString("id INT NOT NULL, \n")
	createTableCommand.WriteString("login varchar(10000) NOT NULL, \n")
	createTableCommand.WriteString("password varchar(10000) NOT NULL \n")
	createTableCommand.WriteString(" );")
	ExecuteDatabaseCommand(db, createTableCommand.String())
}

func InsertIntoUsers(db *sql.DB, id int, login string, password string) {
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, \"%s\", \"%s\")", userTableName, id, login, password)
	ExecuteDatabaseCommand(db, insertCommand)
}

func QueryUsers(db *sql.DB) (int, string, string){
	query := "SELECT * FROM " + userTableName;
	return ExecuteUsersQuery(db, query)
}

// Executes the specified database command
func ExecuteUsersQuery(db *sql.DB, query string) (int, string, string) {
	results, err := db.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var (
		id int
		login string
		password string
	)
	for results.Next(){
		err = results.Scan(&id, &login, &password)
		if err!= nil {
			fmt.Printf("Failed to parse results %s: %s", query, err)
			panic(err)
		}
		fmt.Printf("Id: %d, Login: %s, Password: %s\n", id, login, password)
	}
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}
	return id, login, password
}
