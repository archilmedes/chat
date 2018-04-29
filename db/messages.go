package db

import (
	"fmt"
	"log"
	"database/sql"
)

const (
	Sent     = 0
	Received = 1
)

// Stores a DB Message
type DBMessage struct {
	SSID           uint64
	SentOrReceived int
	Text           []byte
	Timestamp      string
}

// Inserts a message into the messages table
func InsertMessage(SSID uint64, message []byte, timestamp string, sentOrReceived int) bool {
	if sentOrReceived != Sent && sentOrReceived != Received {
		log.Fatalf("Invalid entry for sent/received - must be 0 or 1. Instead, received a %d", sentOrReceived)
	}
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, \"%s\", \"%s\", %d)", messagesTableName, SSID, message, timestamp, sentOrReceived)
	return ExecuteChangeCommand(insertCommand, "Failed to insert into messages")
}

// Removes a message from the messages table
func DeleteMessage(SSID uint64, message string, timestamp string, sentOrReceived int) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE SSID=%d AND message=\"%s\" AND message_timestamp=\"%s\" AND sent_or_received=%d", messagesTableName, SSID, message, timestamp, sentOrReceived)
	return ExecuteChangeCommand(deleteCommand, "Failed to delete message")
}

// Returns all data in the messages table
func QueryMessages() []DBMessage {
	query, err := DB.Prepare("SELECT * FROM messages")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for QueryMessages: %s", err)
	}
	results, err :=query.Query()
	if err != nil {
		fmt.Printf("Error executing QueryMessages query: %s", err)
	}
	return ExecuteMessagesQuery(results)
}

// Returns all messages for a given session identified by SSID
func getSessionMessages(SSID uint64) []DBMessage {
	query, err := DB.Prepare("SELECT * FROM messages WHERE SSID=?")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for getSessionMessages: %s", err)
	}
	results, err :=query.Query(SSID)
	if err != nil {
		fmt.Printf("Error executing getSessionMessages query: %s", err)
	}
	return ExecuteMessagesQuery(results)
}

// Executes the specified database command
func ExecuteMessagesQuery(results *sql.Rows) []DBMessage {
	var messages []DBMessage
	msg := DBMessage{}
	for results.Next() {
		err := results.Scan(&msg.SSID, &msg.Text, &msg.Timestamp, &msg.SentOrReceived)
		if err != nil {
			log.Panicf("Failed to parse results from messages: %s", err)
			panic(err)
		}
		messages = append(messages, msg)
	}
	err := results.Err()
	if err != nil {
		log.Panicf("Failed to get results from conversations query: %s", err)
	}
	return messages
}
