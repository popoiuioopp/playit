package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	var (
		host     = os.Getenv("PostgresqlHost")
		port     = os.Getenv("PostgresqlPort")
		user     = os.Getenv("PostgresqlUser")
		password = os.Getenv("PostgresqlPassword")
		dbname   = os.Getenv("PostgresqlDB")
	)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Database connection established")
}
