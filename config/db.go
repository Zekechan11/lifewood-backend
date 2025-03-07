package config

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB initializes the database connection
func ConnectDB() *sqlx.DB {
	dsn := "root@tcp(localhost:3306)/job_application_db"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connected successfully!")
	return db
}
