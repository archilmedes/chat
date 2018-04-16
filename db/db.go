package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	conf "github.com/wavyllama/chat/config"
	"log"
	"os/exec"
	"time"
)

const (
	databaseName      = "otrmessenger" // Constant in execution, can change
	sessionsTableName = "sessions"
	usersTableName    = "users"
	messagesTableName = "messages"
	testDatabaseName  = "otrmessengertest"
)

var DB *sql.DB

// Function to be called to set everything up
func SetupDatabase() {
	cmd := exec.Command("sh", "db/db_setup.sh", conf.Username, conf.Password)
	SetupDatabaseHelper(databaseName, cmd)
}

func SetupTestDatabase() {
	cmd := exec.Command("sh", "db_test_setup.sh", conf.Username, conf.Password)
	SetupDatabaseHelper(testDatabaseName, cmd)
}

func SetupTestDatabaseAltDir() {
	cmd := exec.Command("sh", "db/db_test_setup.sh", conf.Username, conf.Password)
	SetupDatabaseHelper(testDatabaseName, cmd)
}

func SetupDatabaseHelper(dbName string, cmd *exec.Cmd) {
	err := cmd.Run()
	time.Sleep(2 * time.Second)
	if err != nil {
		fmt.Printf("Error running script: %s", err)

	}
	connectionString := FormConnectionString(dbName)
	DB, _ = ConnectToDatabase(connectionString)
	useDatabaseCommand := "USE " + dbName
	ExecuteChangeCommand(useDatabaseCommand, "Failed to switch databases")
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

func ExecuteChangeCommand(command string, errorMessage string) bool {
	_, err := DB.Exec(command)
	if err != nil {
		fmt.Printf("%s : %s", errorMessage, err)
		return false
	}
	return true
}
