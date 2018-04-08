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

var password = "" // may change for each user
var port = "3306" // may change for each user

// Creates the initial connection to the database
func SetupDatabase() *sql.DB {
	connectionString := FormConnectionString("")

	// Initial connection to MySql - will work even if no databases created
	db, err := ConnectToDatabase(connectionString)

	// Creates the database if it doesn't exist
	createDatabaseString := "CREATE DATABASE IF NOT EXISTS " + databaseName;
	_, err = db.Exec(createDatabaseString)
	if err != nil {
		fmt.Printf("Failed to create db: %s", err)
		panic(err)
	}
	db.Close()

	// Connects to the OTRMessenger database
	connectionString = FormConnectionString(databaseName)
	db, err = ConnectToDatabase(connectionString)

	return db
}

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
