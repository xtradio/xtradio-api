package main

import (
	"database/sql"
	"fmt"
	"log"
)

// UpcomingSongs details
type UpcomingSongs struct {
	Title  string `json:"song"`
	Artist string `json:"artist"`
	Image  string `json:"image"`
	Length int    `json:"length"`
	Share  string `json:"share"`
}

func fetchUpcomingSongsFromDb(list []string) []UpcomingSongs {

	var r []UpcomingSongs

	username := getEnv("MYSQL_USERNAME")

	password := getEnv("MYSQL_PASSWORD")

	host := getEnv("MYSQL_HOST")

	database := getEnv("MYSQL_DATABASE")

	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", username, password, host, database)

	// Open and connect do DB
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Printf("Opening db connection failed: %s", err)
	}

	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Printf("Ping database failed: %s", err)
	}

	for _, v := range list {

		log.Printf("Searching for filename: %s", v)

		var u UpcomingSongs

		query := db.QueryRow("SELECT artist, title, lenght, share, image FROM details WHERE filename=?", v)

		err = query.Scan(&u.Artist, &u.Title, &u.Length, &u.Share, &u.Image)
		if err != nil {
			log.Println("Fetching item failed.", err)
		}

		r = append(r, u)

	}

	return r

}

func upcomingSongs() []UpcomingSongs {
	command := "BACKUP.next"
	data := telnet(command)
	dbParsedData := fetchUpcomingSongsFromDb(data)
	return dbParsedData
}
