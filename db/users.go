package db

import (
	"bytes"
	"fmt"
	"log"
)

type DatabaseUser struct {
	Username, Password, IP string
}

// Creates the users table
func SetupUsersTable() {
	log.Println("Creating the users table...")
	var createTableCommand bytes.Buffer
	createTableCommand.WriteString("CREATE TABLE IF NOT EXISTS ")
	createTableCommand.WriteString(userTableName)
	createTableCommand.WriteString(" (\n")
	createTableCommand.WriteString("Username varchar(1000) NOT NULL, \n")
	createTableCommand.WriteString("Password varchar(1000) NOT NULL, \n")
	createTableCommand.WriteString("ipaddress varchar(18) NOT NULL \n")
	createTableCommand.WriteString(" );")
	ExecuteDatabaseCommand(createTableCommand.String())
}

func UserExists(username string) bool {
	query := "SELECT * FROM " + userTableName + " WHERE Username=\"" + username + "\"";
	users := ExecuteUsersQuery(query)
	return len(users) > 0
}

func GetUser(username string, password string) (*DatabaseUser) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE Username= \"%s\" and Password= \"%s\"", userTableName, username, password);
	users := ExecuteUsersQuery(query)
	if len(users) == 0 {
		return nil
	}
	return &users[0]
}

func AddUser(username string, password string, ipAddress string) bool {
	log.Println("Inserting data into users...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", \"%s\")", userTableName, username, password, ipAddress)
	//ExecuteDatabaseCommand(insertCommand)
	_, err := DB.Exec(insertCommand)
	if err != nil {
		fmt.Printf("Failed to add user %s: %s", username, err)
		return false
	}
	return true
}

func DeleteUser(username string) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE Username= \"%s\"", userTableName, username)
	_, err := DB.Exec(deleteCommand)
	if err != nil {
		fmt.Printf("Failed to delete user %s: %s", username, err)
		return false
	}
	return true
}

func QueryUsers() [] DatabaseUser {
	log.Println("Retrieving data from users...")
	query := "SELECT * FROM " + userTableName;
	return ExecuteUsersQuery(query)
}

// Executes the specified database command
func ExecuteUsersQuery(query string) [] DatabaseUser {
	results, err := DB.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var users [] DatabaseUser
	user := DatabaseUser{}

	for results.Next() {
		err = results.Scan(&user.Username, &user.Password, &user.IP)
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
