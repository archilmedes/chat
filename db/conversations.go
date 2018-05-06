// The db package abstracts out the storage layer and provides CRUD operations
// to various message abstractions
package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/wavyllama/chat/config"
)

// Stores a conversation between two users
type Conversation struct {
	Message Message
	Session Session
}

// Gets all messages between the user and friend
func getConversationsWithFriend(username string, friendDisplayName string) []Conversation {
	query, err := DB.Prepare("SELECT m.*, s.* FROM messages m, sessions s WHERE m.SSID = s.SSID AND m.SSID IN (SELECT s1.SSID FROM sessions s1 WHERE (s1.username=? AND s1.friend_display_name=?) OR (s1.friend_display_name=? AND s1.username=?) ORDER BY s.session_timestamp, m.message_timestamp)")
	if err != nil {
		fmt.Printf("Error creating conversations prepared statement for GetConversationUsers: %s", err)
	}
	results, err := query.Query(username, friendDisplayName, friendDisplayName, username)
	if err != nil {
		fmt.Printf("Error executing GetConversationUsers query: %s", err)
	}
	conversations := ExecuteConversationsQuery(results)
	return conversations
}

// Executes "SELECT" queries for Conversation
func ExecuteConversationsQuery(results *sql.Rows) []Conversation {
	var conversations []Conversation
	convo := Conversation{}
	for results.Next() {
		var timestamp string
		var encMsgText []byte
		err := results.Scan(&convo.Message.SSID, &encMsgText, &timestamp, &convo.Message.SentOrReceived, &convo.Session.SSID, &convo.Session.Username, &convo.Session.FriendDisplayName, &convo.Session.ProtocolType, &convo.Session.ProtocolValue, &convo.Session.timestamp)
		if err != nil {
			log.Panicf("Failed to parse results from conversations:  %s", err)
		}
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", timestamp)
		convo.Message.Timestamp = parsedTime
		convo.Message.Text = AESDecrypt(encMsgText, []byte(db.AESKey))
		conversations = append(conversations, convo)
	}
	err := results.Err()
	if err != nil {
		log.Panicf("Failed to get results from conversations query: %s", err)
	}
	return conversations
}
