package main

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//connecting to db and returns pool of connections
func NewDB() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("USER")+":"+os.Getenv("PASSWORD")+"@tcp(127.0.0.1:"+os.Getenv("PORT")+")/bookshop")
	//db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/bookshop")
	if err != nil {
		panic(err)
	}
	//checking established connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	//Settings for pool
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
