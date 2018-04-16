package db

import (
	"fmt"
	"log"
)

const (
	Sent     = 0
	Received = 1
)

type DBMessage struct {
	SSID, sentOrReceived int
	message, timestamp   string
}

func InsertMessage(SSID int, message string, timestamp string, sentOrReceived int) bool {
	log.Println("Inserting data into messages...")
	if sentOrReceived != Sent && sentOrReceived != Received {
		fmt.Printf("Invalid entry for sent/received - msut be 0 or 1. Instead, received a %d", sentOrReceived)
	}
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, \"%s\", \"%s\", %d)", messagesTableName, SSID, message, timestamp, sentOrReceived)
	return ExecuteChangeCommand(insertCommand, "Failed to insert into messages")
}

func DeleteMessage(SSID int, message string, timestamp string, sentOrReceived int) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE SSID=%d AND message=\"%s\" AND timestamp=\"%s\" AND sent_or_received=%d", messagesTableName, SSID, message, timestamp, sentOrReceived)
	return ExecuteChangeCommand(deleteCommand, "Failed to delete message")
}

func QueryMessages() []DBMessage {
	log.Println("Retrieving data from messages...")
	query := "SELECT * FROM " + messagesTableName
	return ExecuteMessagesQuery(query)
}

// Executes the specified database command
func ExecuteMessagesQuery(query string) []DBMessage {
	results, err := DB.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var messages []DBMessage
	msg := DBMessage{}
	for results.Next() {
		err = results.Scan(&msg.SSID, &msg.message, &msg.timestamp, &msg.sentOrReceived)
		if err != nil {
			fmt.Printf("Failed to parse results %s: %s", query, err)
			panic(err)
		}
		messages = append(messages, msg)
	}
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}
	return messages
}
