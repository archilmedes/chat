package db

import (
	"fmt"
	"log"
)

// Stores a session between two users
type Session struct {
	SSID                                                        uint64
	Username, FriendMac, ProtocolType, ProtocolValue, timestamp string
}

// Inserts data into the sessions table
func InsertIntoSessions(SSID uint64, username string, friendMac string, protocolType string, protocolValue string, timestamp string) bool {
	log.Println("Inserting data into sessions...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, \"%s\", \"%s\", \"%s\", \"%s\", \"%s\")", sessionsTableName, SSID, username, friendMac, protocolType, protocolValue, timestamp)
	return ExecuteChangeCommand(insertCommand, "Failed to insert into sessions")
}

// Deletes a session
func DeleteSession(SSID uint64) bool {
	log.Println("Deleting row from sessions...")
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE SSID =%d", sessionsTableName, SSID)
	return ExecuteChangeCommand(deleteCommand, "Failed to delete session")
}

// Gets all sessions
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
		err = results.Scan(&session.SSID, &session.Username, &session.FriendMac, &session.ProtocolType, &session.ProtocolValue, &session.timestamp)
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
