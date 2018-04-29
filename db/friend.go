package db

import (
	"fmt"
	"log"
	"database/sql"
)

const Self = "me" // Display name for self.

type Friend struct {
	DisplayName, MAC, IP, Username string
}

// Checks if the two users are friends given friend's display name
func areFriends(username, friendDisplayName string) bool {
	query, err := DB.Prepare("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM friends WHERE username=? AND friend_display_name=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for areFriends: %s", err)
	}
	results, err :=query.Query(username, friendDisplayName)
	if err != nil {
		fmt.Printf("Error executing areFriends query: %s", err)
	}
	friends := executeFriendsQuery(results)
	return len(friends) > 0
}

// Add friend to the user's friend's table
func addFriend(username, displayName, macAddress, ipAddress, friendUsername string) bool {
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", \"%s\", \"%s\", \"%s\")", friendsTableName, username, displayName, macAddress, ipAddress, friendUsername)
	return ExecuteChangeCommand(insertCommand, "Failed to add friend")
}

// Delete a friend by their displayName
func deleteFriend(username, displayName string) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE username=\"%s\" AND friend_display_name= \"%s\"", friendsTableName, username, displayName)
	return ExecuteChangeCommand(deleteCommand, "Failed to delete friend")
}

// Updates the user's IP Address
func updateFriendIP(username, macAddress, ipAddress string) bool {
	updateCommand := fmt.Sprintf("UPDATE %s SET friend_ip_address= \"%s\" WHERE username=\"%s\" AND friend_mac_address=\"%s\"", friendsTableName, ipAddress, username, macAddress)
	return ExecuteChangeCommand(updateCommand, "Failed to update friend")
}

// Get all friends
func getFriends(username string) []Friend {
	query, err := DB.Prepare("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM friends WHERE username=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for getFriends: %s", err)
	}
	results, err :=query.Query(username)
	if err != nil {
		fmt.Printf("Error executing getFriends query: %s", err)
	}
	return executeFriendsQuery(results)
}

func getFirstFriend(friends []Friend) *Friend {
	if len(friends) == 0 {
		return nil
	}
	return &friends[0]
}

// Get a friend based on display name
func getFriendByDisplayName(username, friendDisplayName string) *Friend {
	query, err := DB.Prepare("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM friends WHERE username=? AND friend_display_name=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for fetFriendByDisplayName: %s", err)
	}
	results, err :=query.Query(username, friendDisplayName)
	if err != nil {
		fmt.Printf("Error executing getFriendByDisplayName query: %s", err)
	}
	return getFirstFriend(executeFriendsQuery(results))
}

// Get a friend based on username and MAC address
func getFriendByUsernameAndMAC(username, friendUsername, friendMACAddress string) *Friend {
	query, err := DB.Prepare("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM friends WHERE username=? AND friend_username=? AND friend_mac_address=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for getFriendByUsernameAndMAC: %s", err)
	}
	results, err :=query.Query(username, friendUsername, friendMACAddress)
	if err != nil {
		fmt.Printf("Error executing getFriendByUsernameAndMAC query: %s", err)
	}
	return getFirstFriend(executeFriendsQuery(results))
}

// Executes the specified database command
func executeFriendsQuery(results *sql.Rows) []Friend {
	var friends []Friend
	friend := Friend{}
	for results.Next() {
		err := results.Scan(&friend.DisplayName, &friend.MAC, &friend.IP, &friend.Username)
		if err != nil {
			log.Panicf("Failed to parse results from friends:  %s", err)
		}
		friends = append(friends, friend)
	}
	err := results.Err()
	if err != nil {
		log.Panicf("Failed to get results from conversations query: %s", err)
	}
	return friends
}
