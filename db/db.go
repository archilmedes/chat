package db

import (
	_ "github.com/go-sql-driver/mysql"
	"bytes"
	"database/sql"
	"fmt"
)


var hostnameStart = "tcp(127.0.0.1:" // CONSTANT
var hostnameEnd = ")/" //CONSTANT
var username = "root"
var databaseName = "OTRMessenger" // Constant in execution, can change
var tableName = "Messages" // Constant in execution, can change

var password = "" // may change for each user
var port = "3306" // may change for each user

// Function to be called to set everything up
func SetupDatabase() {
	db := InitializeDatabase()
	SetupMessagesTable(db)
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

// Creates the table schema
func SetupMessagesTable(db *sql.DB) {
	var createTableCommand bytes.Buffer
	createTableCommand.WriteString("CREATE TABLE IF NOT EXISTS ")
	createTableCommand.WriteString(tableName)
	createTableCommand.WriteString(" (\n")
	//TODO: Update the db schema
	createTableCommand.WriteString("username char(50) NOT NULL, \n")
	createTableCommand.WriteString("friendName char(50) NOT NULL, \n")
	createTableCommand.WriteString("personalPrivateKey char(50) NOT NULL, \n")
	createTableCommand.WriteString("friendPublicKey char(50) NOT NULL, \n")
	createTableCommand.WriteString("encryptedMessage char(200) NOT NULL \n")
	createTableCommand.WriteString(" );")
	ExecuteDatabaseCommand(db, createTableCommand.String())
}

// Inserts values given into the table
func InsertValues(db *sql.DB, username string, friendsName string, privateKey string, publicKey string, message string){
	insertCommand := fmt.Sprintf("INSERT INTO %s VALUES (\"%s\", \"%s\", \"%s\", \"%s\", \"%s\")", tableName, username, friendsName, privateKey, publicKey, message)
	ExecuteDatabaseCommand(db, insertCommand)
}
