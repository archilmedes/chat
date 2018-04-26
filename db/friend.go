package db

import (
	"fmt"
	"log"
)

const Self = "me" // Display name for self.

type Friend struct {
	DisplayName, MAC, IP, Username string
}

// Checks if the two users are friends given friend's display name
func areFriends(username, friendDisplayName string) bool {
	query := fmt.Sprintf("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM %s WHERE username=\"%s\" AND friend_display_name=\"%s\"", friendsTableName, username, friendDisplayName)
	friends := executeFriendsQuery(query)
	return len(friends) > 0
}

// Add friend to the user's friend's table
func addFriend(username, displayName, macAddress, ipAddress, friendUsername string) bool {
	log.Println("Inserting data into friends...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", \"%s\", \"%s\", \"%s\")", friendsTableName, username, displayName, macAddress, ipAddress, friendUsername)
	return ExecuteChangeCommand(insertCommand, "Failed to add friend")
}

// Delete a friend by their displayName Address
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
	log.Println("Retrieving data from friends...")
	query := fmt.Sprintf("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM %s WHERE username=\"%s\"", friendsTableName, username)
	return executeFriendsQuery(query)
}

func getFirstFriend(friends []Friend) *Friend {
	if len(friends) == 0 {
		return nil
	}
	return &friends[0]
}

// Get a friend based on display name
func getFriendByDisplayName(username, friendDisplayName string) *Friend {
	query := fmt.Sprintf("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM %s WHERE username=\"%s\" AND friend_display_name=\"%s\"", friendsTableName, username, friendDisplayName)
	return getFirstFriend(executeFriendsQuery(query))
}

// Get a friend based on username and MAC address
func getFriendByUsernameAndMAC(username, friendUsername, friendMACAddress string) *Friend {
	query := fmt.Sprintf("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM %s WHERE username=\"%s\" AND friend_username=\"%s\" AND friend_mac_address=\"%s\"", friendsTableName, username, friendUsername, friendMACAddress)
	return getFirstFriend(executeFriendsQuery(query))
}

// Executes the specified database command
func executeFriendsQuery(query string) []Friend {
	results, err := DB.Query(query)
	if err != nil {
		log.Panicf("Failed to execute friend's query: %s", err)
	}
	var friends []Friend
	friend := Friend{}
	for results.Next() {
		err = results.Scan(&friend.DisplayName, &friend.MAC, &friend.IP, &friend.Username)
		if err != nil {
			log.Panicf("Failed to parse results from friends with query: %s;  %s", query, err)
		}
		friends = append(friends, friend)
	}
	err = results.Err()
	if err != nil {
		log.Panicf("Failed to get results from conversations query: %s", err)
	}
	return friends
}
