package db

import (
	"fmt"
	"log"
)

// Check if user has already created an account
func UserExists(username string) bool {
	query := "SELECT username, ipaddress FROM " + usersTableName + " WHERE username=\"" + username + "\""
	users := ExecuteUsersQuery(query)
	return len(users) > 0
}

// Get user from database
func GetUser(username string, password string) *User {
	query := fmt.Sprintf("SELECT username, ipaddress FROM %s WHERE username= \"%s\" and password= \"%s\"", usersTableName, username, password)
	users := ExecuteUsersQuery(query)
	if len(users) == 0 {
		return nil
	}
	return &users[0]
}

// Add new user to database
func AddUser(username string, password string, ipAddress string) bool {
	log.Println("Inserting data into users...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", \"%s\")", usersTableName, username, password, ipAddress)
	return ExecuteChangeCommand(insertCommand, "Failed to add user")
}

// Update the IPv4 address of a user
func UpdateUserIP(username string, ipAddress string) bool {
	log.Println("Updating user's IP...")
	updateCommand := fmt.Sprintf("UPDATE %s SET ipaddress=\"%s\" WHERE username=\"%s\"", usersTableName, ipAddress, username)
	return ExecuteChangeCommand(updateCommand, "Failed to update user's IP")
}

// Update the password of a user
func UpdateUserPassword(username string, password string) bool {
	log.Println("Updating user's password...")
	updateCommand := fmt.Sprintf("UPDATE %s SET password=\"%s\" WHERE username=\"%s\"", usersTableName, password, username)
	return ExecuteChangeCommand(updateCommand, "Failed to update user's password")
}

// Delete a user
func DeleteUser(username string) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE username= \"%s\"", usersTableName, username)
	return ExecuteChangeCommand(deleteCommand, "Failed to delete user")
}

// Get all users
func QueryUsers() []User {
	log.Println("Retrieving data from users...")
	query := "SELECT username, ipaddress FROM " + usersTableName
	return ExecuteUsersQuery(query)
}

// Executes the specified database command
func ExecuteUsersQuery(query string) []User {
	results, err := DB.Query(query)
	if err != nil {
		log.Panicf("Failed to execute %s on conversations table: %s", query, err)
	}
	var users []User
	user := User{}
	for results.Next() {
		err = results.Scan(&user.Username, &user.IP)
		if err != nil {
			log.Panicf("Failed to parse results from conversations with query: %s;  %s", query, err)
		}
		users = append(users, user)
	}
	err = results.Err()
	if err != nil {
		log.Panicf("Failed to get results from users query: %s", err)
	}
	return users
}
