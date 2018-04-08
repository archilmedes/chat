package db

import (
	_ "github.com/go-sql-driver/mysql"
	"bytes"
	"database/sql"
	"fmt"
	"log"
)

const (
	hostnameStart = "tcp(127.0.0.1:" // CONSTANT
	hostnameEnd = ")/" //CONSTANT
	databaseName = "otrmessenger" // Constant in execution, can change
	sessionsTableName = "sessions"
	userTableName = "users"
	conversationTableName = "conversation"
	// tableName = "Messages" // Constant in execution, can change
)

// Function to be called to set everything up
func SetupDatabase() *sql.DB {
	db := InitializeDatabase()
	SetupTables(db)
	SetupSessionsTable(db)
	return db
}
func SetupTables(db *sql.DB) {
	SetupSessionsTable(db)
	SetupUsersTable(db)
	SetupConversationTable(db)

}

// Drops the database if it exists
func DropDatabase(db *sql.DB) {
	dropDatabaseCommand := "DROP DATABASE IF EXISTS " + databaseName;
	ExecuteDatabaseCommand(db, dropDatabaseCommand)
}

// Creates the initial connection to the database
func InitializeDatabase() *sql.DB {
	connectionString := FormConnectionString("")

	// Initial connection to MySql - will work even if no databases created
	db, _ := ConnectToDatabase(connectionString)

	// FOR TESTING ONLY - CLEARS DATABASE EVERY RUN
	DropDatabase(db);

	// Creates the database if it doesn't exist
	createDatabaseCommand := "CREATE DATABASE IF NOT EXISTS " + databaseName;
	ExecuteDatabaseCommand(db, createDatabaseCommand)
	db.Close()

	// Connects to the OTRMessenger database
	connectionString = FormConnectionString(databaseName)
	db, _ = ConnectToDatabase(connectionString)

	useDatabaseCommand := "USE " + databaseName;
	ExecuteDatabaseCommand(db, useDatabaseCommand)

	return db
}

// Executes the specified database command
func ExecuteDatabaseCommand(db *sql.DB, command string){
	_, err := db.Exec(command)
	if err != nil {
		fmt.Printf("Failed to execute command %s: %s", command, err)
		panic(err)
	}
}

// Connects to a database - quits if it encounters errors
func ConnectToDatabase(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Printf("Could not connect to db: %s", err)
		panic(err)
	}
	return db, err
}

// Creates the connection string using username, password, hostname, and port
func FormConnectionString(Name string) string {
	var connectionString bytes.Buffer
	connectionString.WriteString(username)
	connectionString.WriteString(":")
	connectionString.WriteString(password)
	connectionString.WriteString("@")
	connectionString.WriteString(hostnameStart)
	connectionString.WriteString(port)
	connectionString.WriteString(hostnameEnd)
	connectionString.WriteString(Name)
	return connectionString.String()
}

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

// Creates the users table
func SetupUsersTable(db *sql.DB) {
	var createTableCommand bytes.Buffer
	createTableCommand.WriteString("CREATE TABLE IF NOT EXISTS ")
	createTableCommand.WriteString(userTableName)
	createTableCommand.WriteString(" (\n")
	createTableCommand.WriteString("id INT NOT NULL, \n")
	createTableCommand.WriteString("login varchar(10000) NOT NULL, \n")
	createTableCommand.WriteString("password varchar(10000) NOT NULL \n")
	createTableCommand.WriteString(" );")
	ExecuteDatabaseCommand(db, createTableCommand.String())
}

// Creates the conversations table
func SetupConversationTable(db *sql.DB) {
	var createTableCommand bytes.Buffer
	createTableCommand.WriteString("CREATE TABLE IF NOT EXISTS ")
	createTableCommand.WriteString(conversationTableName)
	createTableCommand.WriteString(" (\n")
	createTableCommand.WriteString("SSID INT NOT NULL, \n")
	createTableCommand.WriteString("message varchar(10000) NOT NULL, \n")
	createTableCommand.WriteString("timestamp varchar(30) NOT NULL, \n")
	createTableCommand.WriteString("sent_or_received BIT NOT NULL")
	createTableCommand.WriteString(" );")
	ExecuteDatabaseCommand(db, createTableCommand.String())
}

func InsertIntoConversations(db *sql.DB, SSID int, message string, timestamp string, sentOrReceived int) {
	if sentOrReceived != 0 && sentOrReceived != 1 {
		fmt.Println("Invalid entry for sent/received - msut be 0 or 1. Instead, received a %d", sentOrReceived)
	}
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, \"%s\", \"%s\", %d)", conversationTableName, SSID, message, timestamp, sentOrReceived)
	ExecuteDatabaseCommand(db, insertCommand)
}
func InsertIntoSessions(db *sql.DB, SSID int, userId int, friendId int, privateKey string, fingerprint string) {
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, %d, %d, \"%s\", \"%s\")", sessionsTableName, SSID, userId, friendId, privateKey, fingerprint)
	ExecuteDatabaseCommand(db, insertCommand)
}
func InsertIntoUsers(db *sql.DB, id int, login string, password string) {
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (%d, \"%s\", \"%s\")", userTableName, id, login, password)
	ExecuteDatabaseCommand(db, insertCommand)
}

/*
func QueryDatabase(db *sql.DB) {
	query := "SELECT username, friendName FROM " + tableName;
	ExecuteQuery(db, query)
}
*/

// Executes the specified database command
func ExecuteQuery(db *sql.DB, query string) {
	results, err := db.Query(query)
	if err != nil {
		fmt.Printf("Failed to execute query %s: %s", query, err)
		panic(err)
	}
	var (
		username string
		friendName string
	)
	for results.Next(){
		err = results.Scan(&username, &friendName)
		if err!= nil {
			fmt.Printf("Failed to execute query %s: %s", query, err)
			panic(err)
		}
		fmt.Printf("User: %s, Friend: %s\n", username, friendName)
	}
	err = results.Err()
	if err != nil {
		log.Fatal(err)
	}

}
func ShowTables(db *sql.DB) {
	results, err := db.Query("SHOW TABLES")
	if err != nil {
		fmt.Println("Failed to display tables.")
		panic(err)
	}
	var str string
	for results.Next() {
		err = results.Scan(&str)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Show tables: %s\n", str)
	}
}