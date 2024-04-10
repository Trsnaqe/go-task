package main

import (
	"database/sql"
	"io"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/trsnaqe/gotask/cmd/api"
	"github.com/trsnaqe/gotask/config"
	"github.com/trsnaqe/gotask/db"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()
	wrt := io.MultiWriter(os.Stdout, file)

	log.SetOutput(wrt)

	db, err := db.NewMySQL(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Database connected")
}
