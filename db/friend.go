package db

import (
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
		Logger.Printf("Error creating friends prepared statement for areFriends: %s", err)
	}
	results, err := query.Query(username, friendDisplayName)
	if err != nil {
		Logger.Printf("Error executing areFriends query: %s", err)
	}
	friends := executeFriendsQuery(results)
	return len(friends) > 0
}

// Add friend to the user's friend's table
func addFriend(username, displayName, macAddress, ipAddress, friendUsername string) bool {
	Logger.Println("trying to add friend")
	insertCommand, err := DB.Prepare("INSERT INTO friends VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		Logger.Printf("Error creating friends prepared statement for addFriend: %s", err)
	}
	_, err = insertCommand.Exec(username, displayName, macAddress, ipAddress, friendUsername)
	if err != nil {
		Logger.Panicf("Failed to add friend: %s", err)
	}
	Logger.Println("returning from friend")
	return true
}

// Delete a friend by their displayName
func deleteFriend(username, displayName string) bool {
	deleteCommand, err := DB.Prepare("DELETE FROM friends WHERE username=? AND friend_display_name= ?")
	if err != nil {
		Logger.Printf("Error creating friends prepared statement for deleteFriend: %s", err)
	}
	_, err = deleteCommand.Exec(username, displayName)
	//if err != nil {
	//	Logger.Panicf("Failed to delete friend: %s", err)
	//}
	return true
}

// Updates the user's IP Address
func updateFriendIP(username, macAddress, ipAddress string) bool {
	updateCommand, err := DB.Prepare("UPDATE friends SET friend_ip_address= ? WHERE username=? AND friend_mac_address=?")
	if err != nil {
		Logger.Printf("Error creating friends prepared statement for updateFriendIP: %s", err)
	}
	_, err = updateCommand.Exec(ipAddress, username, macAddress)
	if err != nil {
		Logger.Panicf("Failed to update friend: %s", err)
	}
	return true
}

// Get all friends
func getFriends(username string) []Friend {
	query, err := DB.Prepare("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM friends WHERE username=?")
	if err != nil {
		Logger.Printf("Error creating friends prepared statement for getFriends: %s", err)
	}
	results, err := query.Query(username)
	if err != nil {
		Logger.Printf("Error executing getFriends query: %s", err)
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
		Logger.Printf("Error creating friends prepared statement for fetFriendByDisplayName: %s", err)
	}
	results, err := query.Query(username, friendDisplayName)
	if err != nil {
		Logger.Printf("Error executing getFriendByDisplayName query: %s", err)
	}
	return getFirstFriend(executeFriendsQuery(results))
}

// Get a friend based on username and MAC address
func getFriendByUsernameAndMAC(username, friendUsername, friendMACAddress string) *Friend {
	query, err := DB.Prepare("SELECT friend_display_name, friend_mac_address, friend_ip_address, friend_username FROM friends WHERE username=? AND friend_username=? AND friend_mac_address=?")
	if err != nil {
		Logger.Printf("Error creating friends prepared statement for getFriendByUsernameAndMAC: %s", err)
	}
	results, err := query.Query(username, friendUsername, friendMACAddress)
	if err != nil {
		Logger.Printf("Error executing getFriendByUsernameAndMAC query: %s", err)
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
			Logger.Panicf("Failed to parse results from friends:  %s", err)
		}
		friends = append(friends, friend)
	}
	err := results.Err()
	if err != nil {
		Logger.Panicf("Failed to get results from conversations query: %s", err)
	}
	return friends
}
