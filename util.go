package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func getEnv(envKey string) (envValue string) {

	envValue, ok := os.LookupEnv(envKey)
	if ok != true {
		log.Printf("please set %s environment variable", envKey)
		return
	}

	return envValue
}

func songHistory(list []Song, song Song) []Song {

	var data []Song

	if len(list) == 3 {
		list = list[:len(list)-1]
	}

	data = append(data, song)

	for _, v := range list {
		data = append(data, v)
	}

	return data
}

func dbConnection() (*sql.DB, error) {
	username := getEnv("MYSQL_USERNAME")

	password := getEnv("MYSQL_PASSWORD")

	host := getEnv("MYSQL_HOST")

	database := getEnv("MYSQL_DATABASE")

	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", username, password, host, database)

	// Open and connect do DB
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
