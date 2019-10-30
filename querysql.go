package main

import (
	"database/sql"
	"fmt"
	"os"
)

func querysql(query string) *sql.Rows {

	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))
	// Open and connect do DB
	db, err := sql.Open("mysql", connection)
	if err != nil {
		// http.Error(w, err.Error(), 500)
		fmt.Println("Ping database failed.", err)
		// return "Open database failed"
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		// http.Error(w, err.Error(), 500)
		fmt.Println("Ping database failed.", err)
		// return "Ping database failed"
	}

	rows, err := db.Query(query)
	if err != nil {
		// http.Error(w, err.Error(), 500)
		fmt.Println("Fetching rows failed.", err)
		// return "Fetching rows failed"
	}

	return rows
}
