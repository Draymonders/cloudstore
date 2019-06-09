package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const mysqlLink = "root:123456@tcp(192.168.99.130:13306)/icloud?charset=utf8"

// init : init of mysql
func init() {
	db, _ = sql.Open("mysql", mysqlLink)
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

// ParseRows : TODO parse the return cow
/*
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	return nil
}
*/
// checkErr : if error occurs, then log it
/*
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
*/
