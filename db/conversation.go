package db

import (
	"fmt"
	"log"
)

const (
	Sent     = 0
	Received = 1
)

type Conversation struct {
	SSID, sentOrReceived int
	message, timestamp   string
}

func InsertIntoConversations(SSID int, message string, timestamp string, sentOrReceived int) {
	log.Println("Inserting data into conversations...")
	if sentOrReceived != Sent && sentOrReceived != Received {
		fmt.Printf("Invalid entry for sent/received - msut be 0 or 1. Instead, received a %d", sentOrReceived)
	}
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, \"%s\", \"%s\", %d)", conversationTableName, SSID, message, timestamp, sentOrReceived)
	ExecuteDatabaseCommand(insertCommand)
}

func QueryConversations() []Conversation {
	log.Println("Retrieving data from conversations...")
	query := "SELECT * FROM " + conversationTableName
	return ExecuteConversationQuery(query)
}

// Executes the specified database command
func ExecuteConversationQuery(query string) []Conversation {
	results, err := DB.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var conversations []Conversation
	conv := Conversation{}
	for results.Next() {
		err = results.Scan(&conv.SSID, &conv.message, &conv.timestamp, &conv.sentOrReceived)
		if err != nil {
			fmt.Printf("Failed to parse results %s: %s", query, err)
			panic(err)
		}
		conversations = append(conversations, conv)
	}
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}
	return conversations
}

func GetConversation(userId int, friendId int) []Conversation {
	query := fmt.Sprintf("SELECT * FROM conversations WHERE SSID IN (SELECT SSID FROM sessions WHERE (user_id=%d AND friend_id=%d) OR (friend_id=%d AND user_id=%d))", userId, friendId, userId, friendId)
	conversations := ExecuteConversationQuery(query)
	return conversations
}
