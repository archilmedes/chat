package db

import (
	"bytes"
	"fmt"
	"log"
)

type User struct {
	login string
	displayName string
	password string
}

// Creates the users table
func SetupUsersTable() {
	log.Println("Creating the users table...")
	var createTableCommand bytes.Buffer
	createTableCommand.WriteString("CREATE TABLE IF NOT EXISTS ")
	createTableCommand.WriteString(userTableName)
	createTableCommand.WriteString(" (\n")
	createTableCommand.WriteString("login varchar(1000) NOT NULL, \n")
	createTableCommand.WriteString("displayName varchar(1000) NOT NULL, \n")
	createTableCommand.WriteString("password varchar(1000) NOT NULL \n")
	createTableCommand.WriteString(" );")
	ExecuteDatabaseCommand(createTableCommand.String())
}

func UserExists(username string) bool {
	query := "SELECT COUNT(*) FROM " + userTableName + "WHERE username=" + username;
	users := ExecuteUsersQuery(query)
	return len(users) > 0
}

func InsertIntoUsers(login string, displayName string, password string) {
	log.Println("Inserting data into users...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%s, \"%s\", \"%s\")", userTableName, login, displayName, password)
	ExecuteDatabaseCommand(insertCommand)
}

func QueryUsers() [] User {
	log.Println("Retrieving data from users...")
	query := "SELECT * FROM " + userTableName;
	return ExecuteUsersQuery(query)
}

// Executes the specified database command
func ExecuteUsersQuery(query string) [] User {
	results, err := db.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var users [] User
	user := User{}

	for results.Next() {
		err = results.Scan(&user.displayName, &user.login, &user.password)
		if err!= nil {
			fmt.Printf("Failed to parse results %s: %s", query, err)
			panic(err)
		}
		users = append(users, user)
	}
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}
	return users
}
