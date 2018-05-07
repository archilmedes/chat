package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	conf "github.com/wavyllama/chat/config"
	"log"
	"os/exec"
	"bytes"
)

const (
	databaseName      = "otrmessenger" // Constant in execution, can change
	sessionsTableName = "sessions"
	usersTableName    = "users"
	messagesTableName = "messages"
	friendsTableName  = "friends"
	testDatabaseName  = "otrmessengertest"
	numTables         = 4
)

var DB *sql.DB

// Sets up the database - called at startup
func SetupDatabase() {
	cmd := exec.Command("mysql", fmt.Sprintf("-u%s", conf.Username), fmt.Sprintf("-p%s", conf.Password), "-e", "source scripts/db_setup.sql")
	createDatabase(databaseName, cmd)
}

// Sets up the test database
func SetupTestDatabase() {
	cmd := exec.Command("mysql", fmt.Sprintf("-u%s", conf.Username), fmt.Sprintf("-p%s", conf.Password), "-e", "source ../scripts/db_test_setup.sql")
	createDatabase(testDatabaseName, cmd)
}

// Sets up an empty test database
func SetupEmptyTestDatabase() {
	SetupTestDatabase()
	ClearDatabase()
}

// Runs the command and connects to the database
func createDatabase(dbName string, cmd *exec.Cmd) {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Panicf("%s: %s", err, stderr.String())
	}

	connectionString := formConnectionString(dbName)
	DB, err = connectToDatabase(connectionString)
	if err != nil {
		log.Fatalf("Could not connect to DB: %s", err)
	}
	numTablesCreated := len(ShowTables())
	if numTablesCreated != numTables {
		log.Panicf("Tables were not created properly: expected %d and got %d", numTables, numTablesCreated)
	}
}

// Clears all rows in all tables of the database
func ClearDatabase() {
	ExecuteChangeCommand(fmt.Sprintf("TRUNCATE %s", usersTableName), "Could not truncate table")
	ExecuteChangeCommand(fmt.Sprintf("TRUNCATE %s", messagesTableName), "Could not truncate table")
	ExecuteChangeCommand(fmt.Sprintf("TRUNCATE %s", friendsTableName), "Could not truncate table")
	ExecuteChangeCommand(fmt.Sprintf("TRUNCATE %s", sessionsTableName), "Could not truncate table")
}

// Connects to a database
func connectToDatabase(connectionString string) (*sql.DB, error) {
	return sql.Open("mysql", connectionString)
}

// Creates the connection string using Username, Password, hostname, and port
func formConnectionString(dbName string) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", conf.Username, conf.Password, conf.Port, dbName)
	return connectionString
}

// Shows all tables to ensure DB setup was correct
func ShowTables() []string {
	results, err := DB.Query("SHOW TABLES")
	if err != nil {
		log.Panicf("Failed to display tables: %s", err)
	}
	var tables []string
	var str string
	for results.Next() {
		err = results.Scan(&str)
		tables = append(tables, str)
		if err != nil {
			log.Panicf("Failed to store results: %s", err)
		}
	}
	return tables
}

// Executes Insertions/Updated/Deletes
func ExecuteChangeCommand(command string, errorMessage string) bool {
	_, err := DB.Exec(command)
	if err != nil {
		log.Panicf("Failed to execute change command: %s: %s\n", errorMessage, err)
	}
	return true
}
