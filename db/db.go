package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMySQL(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	return db, nil
}
func InitStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Database connected")
}
