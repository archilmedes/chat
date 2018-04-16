package db

import (
	"fmt"
	"log"
)

type Conversation struct {
	Message DBMessage
	Session Session
}

func GetConversationID(userId int, friendId int) []Conversation {
	query := fmt.Sprintf("SELECT m.*, s.* FROM messages m, sessions s WHERE m.SSID = s.SSID AND m.SSID IN (SELECT s1.SSID FROM sessions s1 WHERE (s1.user_id=%d AND s1.friend_id=%d) OR (s1.friend_id=%d AND s1.user_id=%d))", userId, friendId, userId, friendId)
	conversations := ExecuteConversationsQuery(query)
	return conversations
}

func GetConversationSSID(SSID int) []Conversation {
	query := fmt.Sprintf("SELECT m.*, s.* FROM messages m, sessions s WHERE m.SSID = s.SSID AND m.SSID = %d", SSID)
	conversations := ExecuteConversationsQuery(query)
	return conversations
}

func ExecuteConversationsQuery(query string) []Conversation {
	results, err := DB.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var conversations []Conversation
	convo := Conversation{}
	for results.Next() {
		err = results.Scan(&convo.Message.SSID, &convo.Message.message, &convo.Message.timestamp, &convo.Message.sentOrReceived, &convo.Session.SSID, &convo.Session.UserId, &convo.Session.FriendId, &convo.Session.PrivateKey, &convo.Session.Fingerprint)
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
