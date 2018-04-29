package db

import (
	"fmt"
	"log"
	"database/sql"
)

// Check if user has already created an account
func UserExists(username string) bool {
	query, err := DB.Prepare("SELECT username, ipaddress FROM users WHERE username=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for UserExists: %s", err)
	}
	results, err :=query.Query(username)
	if err != nil {
		fmt.Printf("Error executing UserExists query: %s", err)
	}
	users := ExecuteUsersQuery(results)
	return len(users) > 0
}

// Get user from database
func GetUser(username string, password string) *User {
	query, err := DB.Prepare("SELECT username, ipaddress FROM users WHERE username=? and password=SHA2(?, 256)")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for GetUser: %s", err)
	}
	results, err :=query.Query(username, password)
	if err != nil {
		fmt.Printf("Error executing GetUser query: %s", err)
	}
	users := ExecuteUsersQuery(results)
	if len(users) == 0 {
		return nil
	}
	return &users[0]
}

// Add new user to database
func AddUser(username string, password string, ipAddress string) bool {
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", SHA2(\"%s\", 256), \"%s\")", usersTableName, username, password, ipAddress)
	return ExecuteChangeCommand(insertCommand, "Failed to add user")
}

// Update the IPv4 address of a user
func UpdateUserIP(username string, ipAddress string) bool {
	updateCommand := fmt.Sprintf("UPDATE %s SET ipaddress=\"%s\" WHERE username=\"%s\"", usersTableName, ipAddress, username)
	return ExecuteChangeCommand(updateCommand, "Failed to update user's IP")
}

// Update the password of a user
func UpdateUserPassword(username string, password string) bool {
	updateCommand := fmt.Sprintf("UPDATE %s SET password=SHA2(\"%s\", 256) WHERE username=\"%s\"", usersTableName, password, username)
	return ExecuteChangeCommand(updateCommand, "Failed to update user's password")
}

// Delete a user
func DeleteUser(username string) bool {
	sessionsAndMessages := deleteSessionsWithMessages(username)
	friendsAndUser := deleteUserAndFriends(username)
	return sessionsAndMessages && friendsAndUser
}

// Deletes a user and their friends
func deleteUserAndFriends(username string) bool {
	deleteCommand := fmt.Sprintf("DELETE u, f FROM users u LEFT JOIN friends f ON u.username = f.username WHERE u.username=\"%s\"", username)
	return ExecuteChangeCommand(deleteCommand, "Failed to delete user and friends")
}

// Get all users
func QueryUsers() []User {
	query, err := DB.Prepare("SELECT username, ipaddress FROM users")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for QueryUsers: %s", err)
	}
	results, err :=query.Query()
	if err != nil {
		fmt.Printf("Error executing QueryUsers query: %s", err)
	}
	return ExecuteUsersQuery(results)
}

// Executes the specified database command
func ExecuteUsersQuery(results *sql.Rows) []User {
	var users []User
	user := User{}
	for results.Next() {
		err := results.Scan(&user.Username, &user.IP)
		if err != nil {
			log.Panicf("Failed to parse results from conversations: %s", err)
		}
		users = append(users, user)
	}
	err := results.Err()
	if err != nil {
		log.Panicf("Failed to get results from users query: %s", err)
	}
	return users
}
