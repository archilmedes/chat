package db

import (
	"database/sql"
	"bytes"
	"fmt"
	"log"
)

// Creates the sessions table
func SetupSessionsTable(db *sql.DB) {
	var createTableCommand bytes.Buffer
	createTableCommand.WriteString("CREATE TABLE IF NOT EXISTS ")
	createTableCommand.WriteString(sessionsTableName)
	createTableCommand.WriteString(" (\n")
	createTableCommand.WriteString("SSID INT NOT NULL PRIMARY KEY, \n")
	createTableCommand.WriteString("user_id INT NOT NULL, \n")
	createTableCommand.WriteString("friend_id INT NOT NULL, \n")
	createTableCommand.WriteString("private_key varchar(10000) NOT NULL, \n")
	createTableCommand.WriteString("fingerprint varchar(10000) NOT NULL \n")
	createTableCommand.WriteString(" );")
	ExecuteDatabaseCommand(db, createTableCommand.String())
}

func InsertIntoSessions(db *sql.DB, SSID int, userId int, friendId int, privateKey string, fingerprint string) {
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, %d, %d, \"%s\", \"%s\")", sessionsTableName, SSID, userId, friendId, privateKey, fingerprint)
	ExecuteDatabaseCommand(db, insertCommand)
}

func QuerySessions(db *sql.DB) (int, int, int, string, string){
	query := "SELECT * FROM " + sessionsTableName;
	return ExecuteSessionsQuery(db, query)
}

// Executes the specified database command
func ExecuteSessionsQuery(db *sql.DB, query string) (int, int, int, string, string) {
	results, err := db.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var (
		SSID int
		userId int
		friendId int
		privateKey string
		fingerprint string
	)
	for results.Next(){
		err = results.Scan(&SSID, &userId, &friendId, &privateKey, &fingerprint)
		if err!= nil {
			fmt.Printf("Failed to parse results %s: %s", query, err)
			panic(err)
		}
		fmt.Printf("SSID: %d, User: %d, Friend: %d, Private Key: %s, Fingerprint: %s\n", SSID, userId, friendId, privateKey, fingerprint)
	}
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}
	return SSID, userId, friendId, privateKey, fingerprint
}

