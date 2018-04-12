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
)
var db *sql.DB

// Function to be called to set everything up
func SetupDatabase(){
	db = InitializeDatabase()
	SetupTables()
	SetupSessionsTable()
}
func SetupTables() {
	SetupSessionsTable()
	SetupUsersTable()
	SetupConversationTable()

}

// Drops the database if it exists
func DropDatabase() {
	dropDatabaseCommand := "DROP DATABASE IF EXISTS " + databaseName;
	ExecuteDatabaseCommand(dropDatabaseCommand)
}

// Creates the initial connection to the database
func InitializeDatabase() *sql.DB {
	connectionString := FormConnectionString("")

	// Initial connection to MySql - will work even if no databases created
	db, _ = ConnectToDatabase(connectionString)

	// FOR TESTING ONLY - CLEARS DATABASE EVERY RUN
	DropDatabase();

	// Creates the database if it doesn't exist
	log.Println("Creating database...")
	createDatabaseCommand := "CREATE DATABASE IF NOT EXISTS " + databaseName;
	ExecuteDatabaseCommand(createDatabaseCommand)
	db.Close()

	// Connects to the OTRMessenger database
	connectionString = FormConnectionString(databaseName)
	db, _ = ConnectToDatabase(connectionString)

	log.Println("Switching to OTRMessenger database")
	useDatabaseCommand := "USE " + databaseName;
	ExecuteDatabaseCommand(useDatabaseCommand)

	return db
}

// Executes the specified database command
func ExecuteDatabaseCommand(command string){
	_, err := db.Exec(command)
	if err != nil {
		fmt.Printf("Failed to execute command %s: %s", command, err)
		panic(err)
	}
}

// Connects to a database - quits if it encounters errors
func ConnectToDatabase(connectionString string) (*sql.DB, error) {
	database, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Printf("Could not connect to db: %s", err)
		panic(err)
	}
	return database, err
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

func ShowTables() []string {
	log.Println("Fetching all tables for database")
	results, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Panic("Failed to display tables")
	}
	var tables []string
	var str string
	for results.Next() {
		err = results.Scan(&str)
		tables = append(tables, str)
		if err != nil {
			panic(err)
		}
	}
	return tables
}