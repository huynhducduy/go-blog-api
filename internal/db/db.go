package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"promotion-management-api/internal/config"
)

var db *sql.DB

func OpenConnection() {

	conn, err := sql.Open("mysql", config.GetConfig().DB_USER+":"+config.GetConfig().DB_PASS+"@tcp("+config.GetConfig().DB_HOST+":"+config.GetConfig().DB_PORT+")/"+config.GetConfig().DB_NAME+"?parseTime=true")

	if err != nil {
		log.Fatalf("Cannot open connection, %s", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalf("Cannot ping connection, %s", err)
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	// The sum is the maximum number of concurrent connections
	conn.SetConnMaxLifetime(5 * time.Minute)

	db = conn
}

func GetConnection() *sql.DB {
	return db
}