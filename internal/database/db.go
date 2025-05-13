package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"golang.org/x/net/context"
)

func NewSqlDB(addr string, maxOpenConns int, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("mysql", addr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err

	}

	return db, nil

}
