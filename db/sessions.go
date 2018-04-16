package db

import (
	"fmt"
	"log"
)

type Session struct {
	SSID, UserId, FriendId  int
	PrivateKey, Fingerprint string
}

func InsertIntoSessions(SSID int, userId int, friendId int, privateKey string, fingerprint string) bool {
	log.Println("Inserting data into sessions...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, %d, %d, \"%s\", \"%s\")", sessionsTableName, SSID, userId, friendId, privateKey, fingerprint)
	return ExecuteChangeCommand(insertCommand, "Failed to insert into sessions")
}

func InsertServerData(userId int, friendId int, SSID int) bool {
	log.Println("Inserting server data into sessions...")
	insertCommand := fmt.Sprintf("INSERT INTO %s (user_id, friend_id, SSID) VALUES (%d, %d, %d)", sessionsTableName, userId, friendId, SSID)
	return ExecuteChangeCommand(insertCommand, "Failed to insert server data")
}

func UpdateSessionsOtrData(SSID int, privateKey string, fingerprint string) bool {
	log.Println("Updating sessions with otr data...")
	whereClause := fmt.Sprintf("WHERE SSID = %d", SSID)
	insertCommand := fmt.Sprintf("UPDATE %s SET SSID = %d, private_key = \"%s\", fingerprint = \"%s\" %s", sessionsTableName, SSID, privateKey, fingerprint, whereClause)
	return ExecuteChangeCommand(insertCommand, "Failed to update otr data in sessions")
}

func DeleteSession(SSID int) bool {
	log.Println("Deleting row from sessions...")
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE SSID =%d", sessionsTableName, SSID)
	return ExecuteChangeCommand(deleteCommand, "Failed to delete session")
}

func QuerySessions() []Session {
	log.Println("Retrieving data from sessions...")
	query := "SELECT * FROM " + sessionsTableName
	return ExecuteSessionsQuery(query)
}

// Executes the specified database command
func ExecuteSessionsQuery(query string) []Session {
	results, err := DB.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var sessions []Session
	session := Session{}
	for results.Next() {
		err = results.Scan(&session.SSID, &session.UserId, &session.FriendId, &session.PrivateKey, &session.Fingerprint)
		if err != nil {
			fmt.Printf("Failed to parse results %s: %s", query, err)
			panic(err)
		}
		sessions = append(sessions, session)
	}
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}
	return sessions
}
