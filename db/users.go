package db

import (
	"bytes"
	"fmt"
	"log"
)

type User struct {
	username string
	password string
	ipAddress string
}

// Creates the users table
func SetupUsersTable() {
	log.Println("Creating the users table...")
	var createTableCommand bytes.Buffer
	createTableCommand.WriteString("CREATE TABLE IF NOT EXISTS ")
	createTableCommand.WriteString(userTableName)
	createTableCommand.WriteString(" (\n")
	createTableCommand.WriteString("username varchar(1000) NOT NULL, \n")
	createTableCommand.WriteString("password varchar(1000) NOT NULL, \n")
	createTableCommand.WriteString("ipaddress varchar(18) NOT NULL \n")
	createTableCommand.WriteString(" );")
	ExecuteDatabaseCommand(createTableCommand.String())
}

func UserExists(username string) bool {
	query := "SELECT * FROM " + userTableName + " WHERE username=\"" + username + "\"";
	users := ExecuteUsersQuery(query)
	return len(users) > 0
}

func ValidateCredentials(username string, password string) bool{
	query := fmt.Sprintf("SELECT * FROM %s WHERE username= \"%s\" and password= \"%s\"", userTableName, username, password);
	users := ExecuteUsersQuery(query)
	return len(users) > 0
}

func AddUser(username string, password string, ipAddress string) bool {
	log.Println("Inserting data into users...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", \"%s\")", userTableName, username, password, ipAddress)
	//ExecuteDatabaseCommand(insertCommand)
	_, err := db.Exec(insertCommand)
	if err != nil {
		fmt.Printf("Failed to add user %s: %s", username, err)
		return false
	}
	return true
}

func DeleteUser(username string) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE username= \"%s\"", userTableName, username)
	_, err := db.Exec(deleteCommand)
	if err != nil {
		fmt.Printf("Failed to delete user %s: %s", username, err)
		return false
	}
	return true
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
		err = results.Scan(&user.username, &user.password, &user.ipAddress)
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
