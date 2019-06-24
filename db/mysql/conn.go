package mysql

import (
	"cloudstore/config"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// init : init of mysql
func init() {
	db, _ = sql.Open("mysql", config.MysqlLink)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
	fmt.Println("mysql conn has created")
}

// DBConn : return a mysql db connection
func DBConn() *sql.DB {
	return db
}
