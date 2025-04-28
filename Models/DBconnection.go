package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/joho/godotenv"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	// Retrieve MySQL connection parameters from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open the database connection
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
		return err
	}

	// Ping the database to verify connection
	if err := DB.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
		return err
	}

	fmt.Println("Successfully connected to the database!")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
