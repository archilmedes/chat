package db

import (
	"fmt"
	"log"
)

type Session struct {
	SSID, userId, friendId  int
	privateKey, fingerprint string
}

func InsertIntoSessions(SSID int, userId int, friendId int, privateKey string, fingerprint string) {
	log.Println("Inserting data into conversations...")
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, %d, %d, \"%s\", \"%s\")", sessionsTableName, SSID, userId, friendId, privateKey, fingerprint)
	ExecuteDatabaseCommand(insertCommand)
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
		err = results.Scan(&session.SSID, &session.userId, &session.friendId, &session.privateKey, &session.fingerprint)
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
