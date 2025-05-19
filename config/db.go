package config

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func SetupDB() {
	var err error
	DB, err = sql.Open("mysql", "root:Ten0826817189@tcp(127.0.0.1:3306)/petdatabase")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	log.Println("Database connected successfully")
}
