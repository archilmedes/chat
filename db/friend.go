package db

import (
	"fmt"
	"log"
)

type Friend struct {
	DisplayName, MAC, IP string
}

// Checks if the two users are friends given friend's display name
func AreFriendsName(username string, friendDisplayName string) bool {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=\"%s\" AND friend_display_name=\"%s\"", friendsTableName, username, friendDisplayName)
	friends := ExecuteFriendsQuery(query)
	return len(friends) > 0
}

// Checks if the two users are friends given friend's MAC Address
func AreFriendsMac(username string, friendMacAddress string) bool {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=\"%s\" AND friend_mac_address=\"%s\"", friendsTableName, username, friendMacAddress)
	friends := ExecuteFriendsQuery(query)
	return len(friends) > 0
}

// Get friend from database
func GetFriendByDisplayName(friendDisplayName string) *Friend {
	query := fmt.Sprintf("SELECT friend_display_name, friend_mac_address, friend_ip_address FROM %s WHERE friend_display_name= \"%s\"", friendsTableName, friendDisplayName)
	friends := ExecuteFriendsQuery(query)
	if len(friends) == 0 {
		return nil
	}
	return &friends[0]
}

// Add friend to the user's friend's table
func AddFriend(username string, displayName string, macAddress string, ipAddress string) bool {
	log.Println("Inserting data into friends...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", \"%s\")", friendsTableName, username, displayName, macAddress, ipAddress)
	return ExecuteChangeCommand(insertCommand, "Failed to add friend")
}

// Delete a friend by their MAC Address
func DeleteFriendByMac(macAddress string) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE friend_mac_address= \"%s\"", friendsTableName, macAddress)
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
	query := fmt.Sprintf("SELECT friend_display_name, friend_mac_address, friend_ip_address FROM %s WHERE username=\"%s\"", friendsTableName, username)
	return ExecuteFriendsQuery(query)
}

// Executes the specified database command
func ExecuteFriendsQuery(query string) []Friend {
	results, err := DB.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var friends []Friend
	friend := Friend{}
	for results.Next() {
		err = results.Scan(&friend.DisplayName, &friend.MAC, &friend.IP)
		if err != nil {
			fmt.Printf("Failed to parse results %s: %s", query, err)
			panic(err)
		}
		friends = append(friends, friend)
	}
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}
	return friends
}
