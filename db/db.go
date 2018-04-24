package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	conf "github.com/wavyllama/chat/config"
	"log"
	"os/exec"
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
	cmd := exec.Command("sh", "scripts/db_setup.sh", conf.Username, conf.Password)
	createDatabase(databaseName, cmd)
}

// Sets up the test database
func SetupTestDatabase() {
	cmd := exec.Command("sh", "../scripts/db_test_setup.sh", conf.Username, conf.Password)
	createDatabase(testDatabaseName, cmd)
}

// Sets up an empty test database
func SetupEmptyTestDatabase() {
	cmd := exec.Command("sh", "../scripts/db_test_setup.sh", conf.Username, conf.Password)
	createDatabase(testDatabaseName, cmd)
	clearDatabase()
}

// Sets up the database
func createDatabase(dbName string, cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		log.Panicf("Error running script: %s", err)
	}
	connectionString := formConnectionString(dbName)
	DB, _ = connectToDatabase(connectionString)
	useDatabaseCommand := "USE " + dbName
	ExecuteChangeCommand(useDatabaseCommand, "Failed to switch databases")
	numTablesCreated := len(ShowTables())
	if numTablesCreated != numTables {
		log.Panicf("Tables were not created properly: expected %d and got %d", numTables, numTablesCreated)
	}
}

// Clears all rows in all tables of the database
func clearDatabase() {
	ExecuteChangeCommand(fmt.Sprintf("TRUNCATE %s", usersTableName), "Could not truncate table")
	ExecuteChangeCommand(fmt.Sprintf("TRUNCATE %s", messagesTableName), "Could not truncate table")
	ExecuteChangeCommand(fmt.Sprintf("TRUNCATE %s", friendsTableName), "Could not truncate table")
	ExecuteChangeCommand(fmt.Sprintf("TRUNCATE %s", sessionsTableName), "Could not truncate table")
}

// Connects to a database - quits if it encounters errors
func connectToDatabase(connectionString string) (*sql.DB, error) {
	database, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Panicf("Could not connect to DB: %s", err)
	}
	return database, err
}

// Creates the connection string using Username, Password, hostname, and port
func formConnectionString(Name string) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/", conf.Username, conf.Password, conf.Port)
	return connectionString
}

// Shows all tables to ensure DB setup was correct
func ShowTables() []string {
	log.Println("Fetching all tables for database")
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
		log.Panicf("Failed to execute change command: %s", err)
		return false
	}
	return true
}
