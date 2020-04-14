package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
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

		re := regexp.MustCompile(`\/MUSIC\/.*`)

		parsedFilename := re.FindString(v)

		log.Printf("Searching for filename: %s", parsedFilename)

		var u UpcomingSongs

		query := db.QueryRow("SELECT artist, title, lenght, share, image FROM details WHERE filename=?", parsedFilename)

		err = query.Scan(&u.Artist, &u.Title, &u.Length, &u.Share, &u.Image)
		if err != nil {
			log.Println("Fetching item failed.", err)
		}

		if u.Image == "" {
			u.Image = "default.png"
		}

		u.Image = fmt.Sprintf("https://img.xtradio.org/tracks/%s", u.Image)

		r = append(r, u)

	}

	return r

}

func upcomingSongs() []UpcomingSongs {
	command := "playlist(dot)txt.next"
	data := telnet(command)
	dbParsedData := fetchUpcomingSongsFromDb(data)
	return dbParsedData
}
