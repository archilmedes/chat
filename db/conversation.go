// The db package abstracts out the storage layer and provides CRUD operations
// to various message abstractions
package db

import (
	"fmt"
	"log"
)

// Stores a conversation between two users
type Conversation struct {
	Message DBMessage
	Session Session
}

// Gets all messages between the user and friend
func GetConversationUsers(username string, friendMac string) []Conversation {
	query := fmt.Sprintf("SELECT m.*, s.* FROM messages m, sessions s WHERE m.SSID = s.SSID AND m.SSID IN (SELECT s1.SSID FROM sessions s1 WHERE (s1.username=\"%s\" AND s1.friend_mac=\"%s\") OR (s1.friend_mac=\"%s\" AND s1.username=\"%s\") ORDER BY s.session_timestamp, m.message_timestamp)", username, friendMac, friendMac, username)
	conversations := ExecuteConversationsQuery(query)
	return conversations
}

// Executes "SELECT" queries for Conversation
func ExecuteConversationsQuery(query string) []Conversation {
	results, err := DB.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var conversations []Conversation
	convo := Conversation{}
	for results.Next() {
		err = results.Scan(&convo.Message.SSID, &convo.Message.Text, &convo.Message.Timestamp, &convo.Message.SentOrReceived, &convo.Session.SSID, &convo.Session.Username, &convo.Session.FriendMac, &convo.Session.ProtocolType, &convo.Session.ProtocolValue, &convo.Session.timestamp)
		if err != nil {
			fmt.Printf("Failed to parse results %s: %s", query, err)
			panic(err)
		}
		conversations = append(conversations, convo)
	}
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}
	return conversations
}
