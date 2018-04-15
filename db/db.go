package db

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	conf "chat/config"
	"fmt"
	"log"
	"os/exec"
)

const (
	databaseName = "otrmessenger" // Constant in execution, can change
	sessionsTableName = "sessions"
	userTableName = "users"
	conversationTableName = "conversations"
)
var DB *sql.DB

// Function to be called to set everything up
func SetupDatabase(){
	cmd := exec.Command("sh", "db/db_setup.sh", conf.Username, conf.Password)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running script: %s", err)

	}
	connectionString := FormConnectionString(databaseName)
	DB, _ = ConnectToDatabase(connectionString)
	useDatabaseCommand := "USE " + databaseName
	ExecuteDatabaseCommand(useDatabaseCommand)
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
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/", conf.Username, conf.Password, conf.Port)
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