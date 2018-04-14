package db

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	config "chat/config"
	"fmt"
	"log"
)

const (
	databaseName = "otrmessenger" // Constant in execution, can change
	sessionsTableName = "sessions"
	userTableName = "users"
	conversationTableName = "conversation"
)
var DB *sql.DB

// Function to be called to set everything up
func SetupDatabase(){
	DB = InitializeDatabase()
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
	DB, _ = ConnectToDatabase(connectionString)

	// FOR TESTING ONLY - CLEARS DATABASE EVERY RUN
	DropDatabase();

	// Creates the database if it doesn't exist
	log.Println("Creating database...")
	createDatabaseCommand := "CREATE DATABASE IF NOT EXISTS " + databaseName
	ExecuteDatabaseCommand(createDatabaseCommand)
	DB.Close()

	// Connects to the OTRMessenger database
	connectionString = FormConnectionString(databaseName)
	DB, _ = ConnectToDatabase(connectionString)

	log.Println("Switching to OTRMessenger database")
	useDatabaseCommand := "USE " + databaseName
	ExecuteDatabaseCommand(useDatabaseCommand)

	return DB
}

// Executes the specified database command
func ExecuteDatabaseCommand(command string){
	_, err := DB.Exec(command)
	if err != nil {
		fmt.Printf("Failed to execute command %s: %s", command, err)
		panic(err)
	}
}

// Connects to a database - quits if it encounters errors
func ConnectToDatabase(connectionString string) (*sql.DB, error) {
	database, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Printf("Could not connect to DB: %s", err)
		panic(err)
	}
	return database, err
}

// Creates the connection string using Username, Password, hostname, and port
func FormConnectionString(Name string) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/", config.Username, config.Password, config.Port)
	return connectionString
}

func ShowTables() []string {
	log.Println("Fetching all tables for database")
	results, err := DB.Query("SHOW TABLES")
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