package db

import (
	"fmt"
	"log"
	"database/sql"
)

func generateSalt(username string, password string) string {
	saltAndPassword := fmt.Sprintf("%d%s%s%d", len(password), username, password, len(username))
	return saltAndPassword
}

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
	hashedPassword := generateSalt(username, password)
	query, err := DB.Prepare("SELECT username, ipaddress FROM users WHERE username=? and password=SHA2(?,256)")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for GetUser: %s", err)
	}
	results, err :=query.Query(username, hashedPassword)
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
	hashedPassword := generateSalt(username, password)
	insertCommand, err := DB.Prepare("INSERT INTO users VALUES (?, SHA2(?,256), ?)")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for AddUser: %s", err)
	}
	_, err = insertCommand.Exec(username, hashedPassword, ipAddress)
	if err != nil {
		log.Panicf("Failed to add user: %s", err)
	}
	return true
}

// Update the IPv4 address of a user
func UpdateUserIP(username string, ipAddress string) bool {
	updateCommand, err := DB.Prepare("UPDATE users SET ipaddress=? WHERE username=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for UpdateUserIP: %s", err)
	}
	_, err = updateCommand.Exec(ipAddress, username)
	if err != nil {
		log.Panicf("Failed to update user's IP: %s", err)
	}
	return true
}

// Update the password of a user
func UpdateUserPassword(username string, password string) bool {
	hashedPassword := generateSalt(username, password)
	updateCommand, err := DB.Prepare("UPDATE users SET password=SHA2(?,256) WHERE username=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for UpdateUserPassword: %s", err)
	}
	_, err = updateCommand.Exec(hashedPassword, username)
	if err != nil {
		log.Panicf("Failed to update user's password: %s", err)
	}
	return true
}

// Delete a user
func DeleteUser(username string) bool {
	sessionsAndMessages := deleteSessionsWithMessages(username)
	friendsAndUser := deleteUserAndFriends(username)
	return sessionsAndMessages && friendsAndUser
}

// Deletes a user and their friends
func deleteUserAndFriends(username string) bool {
	deleteCommand, err := DB.Prepare("DELETE u, f FROM users u LEFT JOIN friends f ON u.username = f.username WHERE u.username=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for deleteUserAndFriends: %s", err)
	}
	_, err = deleteCommand.Exec(username)
	if err != nil {
		log.Panicf("Failed to delete user and friends: %s", err)
	}
	return true
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
