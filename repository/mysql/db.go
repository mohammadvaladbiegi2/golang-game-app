package mysql

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	db *sql.DB
}

// docker := "gameapp:password@(localhost:3308)/gameapp_db"

func NewDB() *MySQLDB {
	db, err := sql.Open("mysql", "root:@(localhost:3306)/gameapp_db?parseTime=true")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{db}
}
