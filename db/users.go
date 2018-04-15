package db

import (
	"fmt"
	"log"
)

func UserExists(username string) bool {
	query := "SELECT username, ipaddress FROM " + userTableName + " WHERE username=\"" + username + "\"";
	users := ExecuteUsersQuery(query)
	return len(users) > 0
}

func GetUser(username string, password string) *User {
	query := fmt.Sprintf("SELECT username, ipaddress FROM %s WHERE username= \"%s\" and password= \"%s\"", userTableName, username, password);
	users := ExecuteUsersQuery(query)
	if len(users) == 0 {
		return nil
	}
	return &users[0]
}

func AddUser(username string, password string, ipAddress string) bool {
	log.Println("Inserting data into users...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", \"%s\")", userTableName, username, password, ipAddress)
	// ExecuteDatabaseCommand(insertCommand)
	_, err := DB.Exec(insertCommand)
	if err != nil {
		fmt.Printf("Failed to add user %s: %s", username, err)
		return false
	}
	return true
}

func DeleteUser(username string) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE username= \"%s\"", userTableName, username)
	_, err := DB.Exec(deleteCommand)
	if err != nil {
		fmt.Printf("Failed to delete user %s: %s", username, err)
		return false
	}
	return true
}

func QueryUsers() [] User {
	log.Println("Retrieving data from users...")
	query := "SELECT username, ipaddress FROM " + userTableName;
	return ExecuteUsersQuery(query)
}

// Executes the specified database command
func ExecuteUsersQuery(query string) [] User {
	results, err := DB.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var users [] User
	user := User{}

	for results.Next() {
		err = results.Scan(&user.Username, &user.IP)
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
