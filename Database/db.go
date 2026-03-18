package db

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var DB *sqlx.DB // exported global DB connection

func Init() error {
	// Load .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	// Change this line:
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=skip-verify", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}

	// Optional: test connection
	if err := DB.Ping(); err != nil {
		return err
	}

	fmt.Println("Database connected successfully")
	return nil
}
