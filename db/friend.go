package db

import (
	"fmt"
	"log"
)

type Friend struct {
	DisplayName, MAC, IP, Username string
}

// Checks if the two users are friends given friend's display name
func AreFriends(username string, friendDisplayName string) bool {
	query := fmt.Sprintf("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM %s WHERE username=\"%s\" AND friend_display_name=\"%s\"", friendsTableName, username, friendDisplayName)
	friends := ExecuteFriendsQuery(query)
	return len(friends) > 0
}

// Add friend to the user's friend's table
func AddFriend(username string, displayName string, macAddress string, ipAddress string, friendUsername string) bool {
	log.Println("Inserting data into friends...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", \"%s\", \"%s\", \"%s\")", friendsTableName, username, displayName, macAddress, ipAddress, friendUsername)
	return ExecuteChangeCommand(insertCommand, "Failed to add friend")
}

// Delete a friend by their MAC Address
func DeleteFriend(username string, macAddress string) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE username=\"%s\" AND friend_mac_address= \"%s\"", friendsTableName, username, macAddress)
	return ExecuteChangeCommand(deleteCommand, "Failed to delete friend")
}

// Updates the user's IP Address
func UpdateFriendIp(macAddress string, ipAddress string) bool {
	updateCommand := fmt.Sprintf("UPDATE %s SET friend_ip_address= \"%s\" WHERE friend_mac_address=\"%s\"", friendsTableName, ipAddress, macAddress)
	return ExecuteChangeCommand(updateCommand, "Failed to update friend")
}

// Get all friends
func GetFriends(username string) []Friend {
	log.Println("Retrieving data from friends...")
	query := fmt.Sprintf("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM %s WHERE username=\"%s\"", friendsTableName, username)
	return ExecuteFriendsQuery(query)
}

// Executes the specified database command
func ExecuteFriendsQuery(query string) []Friend {
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
