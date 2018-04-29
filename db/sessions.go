package db

import (
	"encoding/hex"
	"fmt"
	"log"
	"database/sql"
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
	query, err := DB.Prepare("SELECT * FROM sessions")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for QuerySessions: %s", err)
	}
	results, err :=query.Query()
	if err != nil {
		fmt.Printf("Error executing QuerySessions query: %s", err)
	}
	return ExecuteSessionsQuery(results)
}

// Deletes the sessions and messages of the given user
func deleteSessionsWithMessages(username string) bool {
	deleteCommand := fmt.Sprintf("DELETE s, m FROM sessions s LEFT JOIN messages m ON s.SSID = m.SSID WHERE s.Username=\"%s\"", username)
	return ExecuteChangeCommand(deleteCommand, "Failed to do large delete")
}

// Get all sessions belonging to a user by the username
func getUserSessions(username string) []Session {
	query, err := DB.Prepare("SELECT * FROM sessions WHERE username=? ORDER BY session_timestamp DESC")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for getUserSessions: %s", err)
	}
	results, err :=query.Query(username)
	if err != nil {
		fmt.Printf("Error executing getUserSessions query: %s", err)
	}
	return ExecuteSessionsQuery(results)
}

// Get the session corresponding to the session identifier
func GetSession(SSID uint64) *Session {
	query, err := DB.Prepare("SELECT * FROM sessions WHERE SSID=%d")
	if err != nil {
		fmt.Printf("Error creating users prepared statement for GetSession: %s", err)
	}
	results, err :=query.Query(SSID)
	if err != nil {
		fmt.Printf("Error executing GetSession query: %s", err)
	}
	sessions := ExecuteSessionsQuery(results)
	if len(sessions) == 0 {
		return nil
	}
	return &sessions[0]
}

// Executes the specified database command
func ExecuteSessionsQuery(results *sql.Rows) []Session {
	var sessions []Session
	session := Session{}
	for results.Next() {
		err := results.Scan(&session.SSID, &session.Username, &session.FriendDisplayName, &session.ProtocolType, &session.ProtocolValue, &session.timestamp)
		if err != nil {
			log.Panicf("Failed to parse results from conversations: %s", err)
		}
		sessions = append(sessions, session)
	}
	err := results.Err()
	if err != nil {
		log.Panicf("Failed to get results from sessions query: %s", err)
	}
	return sessions
}
