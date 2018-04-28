package db

import (
	"encoding/hex"
	"fmt"
	"log"
)

// Stores a session between two users
type Session struct {
	SSID                                                 uint64
	Username, FriendDisplayName, ProtocolType, timestamp string
	ProtocolValue                                        []byte
}

// Inserts data into the sessions table
func InsertIntoSessions(SSID uint64, username string, friendMac string, protocolType string, protocolValue []byte, timestamp string) bool {
	hexProtoValue := hex.EncodeToString(protocolValue)
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, \"%s\", \"%s\", \"%s\", UNHEX(\"%s\"), \"%s\")", sessionsTableName, SSID, username, friendMac, protocolType, hexProtoValue, timestamp)
	return ExecuteChangeCommand(insertCommand, "Failed to insert into sessions")
}

// Deletes a session
func DeleteSession(SSID uint64) bool {
	deleteCommand := fmt.Sprintf("DELETE FROM %s WHERE SSID =%d", sessionsTableName, SSID)
	return ExecuteChangeCommand(deleteCommand, "Failed to delete session")
}

// Gets all sessions
func QuerySessions() []Session {
	query := "SELECT * FROM " + sessionsTableName
	return ExecuteSessionsQuery(query)
}

// Deletes the sessions and messages of the given user
func deleteSessionsWithMessages(username string) bool {
	deleteCommand := fmt.Sprintf("DELETE s, m FROM sessions s LEFT JOIN messages m ON s.SSID = m.SSID WHERE s.Username=\"%s\"", username)
	return ExecuteChangeCommand(deleteCommand, "Failed to do large delete")
}

// Get all sessions belonging to a user by the username
func getUserSessions(username string) []Session {
	queryCommand := fmt.Sprintf("SELECT * FROM %s WHERE Username=\"%s\" ORDER BY sessions.session_timestamp DESC", sessionsTableName, username)
	return ExecuteSessionsQuery(queryCommand)
}

// Executes the specified database command
func ExecuteSessionsQuery(query string) []Session {
	results, err := DB.Query(query)
	if err != nil {
		log.Panicf("Failed to execute %s on conversations table: %s", query, err)
	}
	var sessions []Session
	session := Session{}
	for results.Next() {
		err = results.Scan(&session.SSID, &session.Username, &session.FriendDisplayName, &session.ProtocolType, &session.ProtocolValue, &session.timestamp)
		if err != nil {
			log.Panicf("Failed to parse results from conversations with query: %s;  %s", query, err)
		}
		sessions = append(sessions, session)
	}
	err = results.Err()
	if err != nil {
		log.Panicf("Failed to get results from sessions query: %s", err)
	}
	return sessions
}
